 
/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package tasks

import (
	"context"
	"sync"
	"time"
	"reflect"
	
	"go.uber.org/atomic"

	`github.com/spkaeros/rscgo/pkg/game/entity`
	`github.com/spkaeros/rscgo/pkg/log`
)

// tickable script procedure
type (
	//scriptCall This is a type-free box type, made for boxing function pointers usually from an Anko script.
	ScriptCall interface{}
	//scriptArg This is a type-free box type, made for boxing any function arguments, for use in conjunction with scriptCall.
	scriptArg interface{}
	//scriptArg A slice of ScriptCalls.
	ScriptCalls []ScriptCall
	//Scripts A locked slice of ScriptCalls.
	Scripts struct {
		ScriptCalls
		sync.RWMutex
	}
)

// Call function type aliases for tickables etc
type (
	call = func()
	StatusReturnCall = func() bool
	dualReturnCall = func(context.Context) (reflect.Value,reflect.Value)
	singleArgDualReturnCall = func(context.Context, reflect.Value) (reflect.Value,reflect.Value)
	playerArgCall = func(player entity.MobileEntity)
	playerArgStatusReturnCall = func(player entity.MobileEntity) bool
)

var (
	Ticks = atomic.NewUint64(0)
)

type tickCount int

func CurrentTick() tickCount {
	return tickCount(Ticks.Load())
}

func (t tickCount) since(t1 tickCount) tickCount {
	return t-t1
}

//TickList A collection of Tasks that are intended to be ran once per game engine tick.
// Tasks should contractually return either true if they are to be removed after execution completes,
// or false if they are to be ran again on the next engine cycle.
var TickList = &Scripts{}

func Schedule(ticks int, call StatusReturnCall) {
	TickList.Schedule(ticks, call)
}

//DoOnceSync Run a task one time, and then remove it, regardless of the return value of the call.
func DoOnceSync(ticks int, fn ScriptCall) {
	c := make(chan struct{})
	start := CurrentTick()
	var ticker = func() bool {
		if CurrentTick().since(start) >= tickCount(ticks)  {
			defer close(c)
			switch fn.(type) {
			case call:
				fn.(call)()
			case StatusReturnCall:
				fn.(StatusReturnCall)()
			case dualReturnCall:
				fn.(dualReturnCall)(context.Background())
			case playerArgCall:
				fn.(playerArgCall)(nil)
			case playerArgStatusReturnCall:
				fn.(playerArgStatusReturnCall)(nil)
			default:
				log.Debugf("Couldn't run task[%v]: Type '%T' not handled.", fn, fn)
			}
			// start = CurrentTick()
			return true
		}
		return false
	}
	TickList.Lock()
	TickList.ScriptCalls = append(TickList.ScriptCalls, ticker)
	TickList.Unlock()
	<-c
}

func DoOnce(ticks int, fn ScriptCall) {
	start := CurrentTick()
	var ticker = func() bool {
		if CurrentTick().since(start) >= tickCount(ticks)  {
			switch fn.(type) {
			case call:
				fn.(call)()
			case StatusReturnCall:
				fn.(StatusReturnCall)()
			case dualReturnCall:
				fn.(dualReturnCall)(context.Background())
			case playerArgCall:
				fn.(playerArgCall)(nil)
			case playerArgStatusReturnCall:
				fn.(playerArgStatusReturnCall)(nil)
			default:
				log.Debugf("Couldn't run task[%v]: Type '%T' not handled.", fn, fn)
			}
			// start = CurrentTick()
			return true
		}
		return false
	}
	TickList.Lock()
	TickList.ScriptCalls = append(TickList.ScriptCalls, ticker)
	TickList.Unlock()
}

func Stall(ticks int) {
	c := make(chan struct{})
	startTick := CurrentTick()
	ticker := func() bool {
		if CurrentTick().since(startTick) >= tickCount(ticks)  {
			close(c)
			return true
		}
		return false
	}
	TickList.Lock()
	TickList.ScriptCalls = append(TickList.ScriptCalls, ticker)
	TickList.Unlock()
	<-c
}

func (s *Scripts) Schedule(ticks int, fn ScriptCall) {
	// startTick := CurrentTick()
	// var ticker func() bool
	// ticker = func() bool {
		// if CurrentTick().since(startTick) >= tickCount(ticks)  {
			// startTick = CurrentTick()
			// return fn()
		// }
		// return false
	// }
	// s.Lock()
	// s.ScriptCalls = append(s.ScriptCalls, ticker)
	// s.Unlock()

	start := CurrentTick()
	s.Lock()
	s.ScriptCalls = append(s.ScriptCalls, func() bool {
		if CurrentTick().since(start) >= tickCount(ticks)  {
			start = CurrentTick()
			switch fn.(type) {
			case call:
				fn.(call)()
			case playerArgCall:
				fn.(playerArgCall)(nil)
			case StatusReturnCall:
				return fn.(StatusReturnCall)()
			case dualReturnCall:
				ret, err := fn.(dualReturnCall)(context.Background())
				return !err.IsNil() || ret.Bool()
			case playerArgStatusReturnCall:
				return fn.(playerArgStatusReturnCall)(nil)
			default:
				log.Debugf("Couldn't run task[%v]: Type '%T' not handled.", fn, fn)
			}
			return true
		}
		return false
	})
	s.Unlock()

}

func (s *Scripts) Add(fn interface{}) {
	s.Lock()
	// log.Debugf("%v (%T)\n", fn, fn)
	s.ScriptCalls = append(s.ScriptCalls, fn)
	s.Unlock()
}

func Do(fn interface{}) {
	TickList.Add(fn)
}


func (s *Scripts) Tick(ctx context.Context) {
	// s.ForEach(nil)
	s.ForEach(ctx, nil)
}

type sigFin chan struct{}

func (s *Scripts) ForEach(ctx context.Context, arg interface{}) {
	var list = s.ScriptCalls[:0]
	tickCtx, cancel := context.WithTimeout(ctx, 640*time.Millisecond)
	done := make(sigFin)
	s.RLock()
	// log.Debug(tickCtx.Value("server"))
	go func(ctx context.Context) {
		defer cancel()
		for _, script := range s.ScriptCalls {
			select {
			case <-ctx.Done():
				log.Debug("Task scheduling context reached timeout with the error value:", ctx.Err())
				return
			default:
				if script == nil {
					continue
				}
				list = append(list, script)
				switch script.(type) {
				// Simple function call, no input no input
				case call:
					script.(call)()
				// A function call taking a MobileEntity interface as an argument
				case playerArgCall:
					script.(playerArgCall)(arg.(entity.MobileEntity))
				// A function call returning its active status.
				case StatusReturnCall: 
					if script.(StatusReturnCall)() {
						list = list[:len(list)-1]
					}
				// A function call taking a *world.Player as an argument and returning its active status.
				case playerArgStatusReturnCall:
					if script.(playerArgStatusReturnCall)(arg.(entity.MobileEntity)) {
						list = list[:len(list)-1]
					}
				// A function call that returns two values, the first a result value, and the second an error value
				// Upon non-nil error value, it will log the stringified err struct then remove from active list,
				// otherwise schedules the same call to run again next tick.
				case dualReturnCall:
					ret, callErr := (script.(dualReturnCall))(context.Background())
					if !callErr.IsNil() {
						log.Warn("Error retVal from a dualReturnCall in the Anko ctx:", callErr.Elem())
						continue
					}
					if ret.Bool() {
						list = list[:len(list)-1]
					}
				// A function call that returns two values, the first a result value, and the second an error value
				// Requires one argument, no type restrictions, so long as the client reads it properly.
				// Upon non-nil error value, it will log the stringified err struct then remove from active list,
				// otherwise schedules the same call to run again next tick.
				case singleArgDualReturnCall:
					ret, callErr := (script.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(arg))
					if !callErr.IsNil() {
						log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
						continue
					}
					if ret.Bool() {
						list = list[:len(list)-1]
					}
				default:
					log.Debugf("Couldn't run task[%v]: Type '%T' not handled.", script, script)
				}
		}
	}
		close(done)
	}(tickCtx)
	s.RUnlock()
	<-done
	s.Lock()
	defer s.Unlock()
	s.ScriptCalls = list
	Ticks.Inc()
	// wait.Wait()
}
