package world

import (
	"context"
	"reflect"
	"sync"
	
	"github.com/spkaeros/rscgo/pkg/log"
)

// tickable script procedure
type (
	//scriptCall This is a type-free box type, made for boxing function pointers usually from an Anko script.
	scriptCall interface{}
	//scriptArg This is a type-free box type, made for boxing any function arguments, for use in conjunction with scriptCall.
	scriptArg interface{}
	//scriptArg A slice of scriptCalls.
	scriptCalls []scriptCall
	//scripts A locked slice of scriptCalls.
	scripts struct {
		scriptCalls
		sync.RWMutex
	}
)

// Call function type aliases for tickables etc
type (
	call = func()
	statusReturnCall = func() bool
	dualReturnCall = func(context.Context) (reflect.Value,reflect.Value)
	singleArgDualReturnCall = func(context.Context, reflect.Value) (reflect.Value,reflect.Value)
	playerArgCall = func(player *Player)
	playerArgStatusReturnCall = func(player *Player) bool
)

func (s *scripts) Add(task scriptCall) {
	s.Lock()
	defer s.Unlock()
	s.scriptCalls = append(s.scriptCalls, task)
}

func (s *scripts) Tick(arg interface{}) {
	s.async(arg)
}

func (s *scripts) Schedule(ticks int, fn func() bool) {
	remainder := ticks
	ticker := func() bool {
		remainder -= 1
		if remainder <= 0 {
			remainder = ticks
			return fn()
		}
		return false
	}
	s.Add(ticker)
}

func (s *scripts) async(arg interface{}) {
	retainedScripts := make(chan scriptCall, len(s.scriptCalls))
	wait := sync.WaitGroup{}
	s.RLock()
	for _, script := range s.scriptCalls {
		wait.Add(1)
		go func(script scriptCall) {
			defer wait.Done()
			// Determine the type of our script callback
			switch script.(type) {
			// Simple function call, no input no input
			case call:
				(script.(call))()
				retainedScripts <- script
			// A function call taking a *world.Player as an argument
			case playerArgCall:
				if p, ok := arg.(*Player); ok {
					(script.(playerArgCall))(p)
					retainedScripts <- script
				}
			// A function call returning its active status.
			case statusReturnCall:
				if !(script.(statusReturnCall))() {
					retainedScripts <- script
				}
			// A function call taking a *world.Player as an argument and returning its active status.
			case playerArgStatusReturnCall:
				if p, ok := arg.(*Player); ok && !(script.(playerArgStatusReturnCall))(p) {
					retainedScripts <- script
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case dualReturnCall:
				ret, callErr := (script.(dualReturnCall))(context.Background())
				if !callErr.IsNil() {
					log.Warn("Error retVal from a dualReturnCall in the Anko ctx:", callErr.Elem())
					return
				}
				if v, ok := ret.Interface().(bool); ok && !v {
					log.Debugf("%v, %v\n", v, callErr)
					return
				}
				retainedScripts <- script

//				retainedScripts = append(retainedScripts, script)
			// A function call that returns two values, the first a result value, and the second an error value
			// Requires one argument, no type restrictions, so long as the client reads it properly.
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case singleArgDualReturnCall:
				ret, callErr := (script.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(arg))
				if !callErr.IsNil() {
					log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
					return
				}
				if !ret.IsNil() && ret.Bool() {
					log.Debug(ret)
				}
				retainedScripts <- script
			default: return
			}
	//	}()
		}(script)
	}
	s.RUnlock()
	wait.Wait()
	s.Lock()
	defer s.Unlock()
	s.scriptCalls = s.scriptCalls[:0]
	select {
	case fn, ok := <-retainedScripts:
		if !ok || fn == nil {
			return
		}
		s.scriptCalls = append(s.scriptCalls, fn)
	default: return
	}
}
