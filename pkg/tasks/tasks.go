 
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
	"reflect"
	
	`github.com/spkaeros/rscgo/pkg/game/entity`
	`github.com/spkaeros/rscgo/pkg/log`
)

// tickable script procedure
type (
	//scriptCall This is a type-free box type, made for boxing function pointers usually from an Anko script.
	ScriptCall struct {
		fn interface{}
	}
	//scriptArg This is a type-free box type, made for boxing any function arguments, for use in conjunction with scriptCall.
	scriptArg interface{}
	//scriptArg A slice of ScriptCalls.
	ScriptCalls = []*ScriptCall
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

//Task is a single func that takes no args and returns a bool to indicate whether or not it
// should be removed from the set it belongs to upon completion
type Task func() bool

//Tasks is a slice of Task closures
type Tasks []Task

//TaskList is a concurrency-safe list of Task closures mapped to simple string identifiers.
// The main purpose for this is to provide the ability to schedule functions to run on the
// game engines thread, to synchronize certain sensitive events.
// It is guarded by a sync.RWMutex.
type TaskList struct {
	Tasks
	sync.RWMutex
}

//TickList A collection of Tasks that are intended to be ran once per game engine tick.
// Tasks should contractually return either true if they are to be removed after execution completes,
// or false if they are to be ran again on the next engine cycle.
var TickList = &Scripts{}

func Schedule(ticks int, call Task) {
	TickList.Schedule(ticks, call)
}

func (s *Scripts) Schedule(ticks int, fn Task) {
	curTick := 0
	ticker := func() bool {
		if curTick >= ticks {
			curTick = 0
			return fn()
		}
		defer func() {
			curTick += 1
		}()
		return false
	}
	s.Add(ticker)
}

func (s *Scripts) Add(fn interface{}) {
	s.Lock()
	defer s.Unlock()
	// log.Debugf("%v (%T)\n", fn, fn)
	s.ScriptCalls = append(s.ScriptCalls, &ScriptCall{fn})
}

func (s *Scripts) Tick() {
	s.ForEach(nil)
}

func (s *Scripts) ForEach(arg interface{}) {
	s.RLock()
	// removeList := make(ScriptCalls, 0, len(s.ScriptCalls))
	keepList := make(ScriptCalls, 0, len(s.ScriptCalls))
	wait := sync.WaitGroup{}
	for _, script := range s.ScriptCalls {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if script == nil {
				return
			}
			// Determine the type of our script callback
			keepList = append(keepList, script)
			switch script.fn.(type) {
			// Simple function call, no input no input
			case call:
				(script.fn.(call))()
			// A function call taking a *world.Player as an argument
			case playerArgCall:
				(script.fn.(playerArgCall))(arg.(entity.MobileEntity))
			// A function call returning its active status.
			case StatusReturnCall:
				if (script.fn.(StatusReturnCall))() {
					keepList = keepList[:len(keepList)-1]
				}
			// A function call taking a *world.Player as an argument and returning its active status.
			case playerArgStatusReturnCall:
				if (script.fn.(playerArgStatusReturnCall))(arg.(entity.MobileEntity)) {
					keepList = keepList[:len(keepList)-1]
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case dualReturnCall:
				ret, callErr := (script.fn.(dualReturnCall))(context.Background())
				if !callErr.IsNil() {
					log.Warn("Error retVal from a dualReturnCall in the Anko ctx:", callErr.Elem())
					return
				}
				if v, ok := ret.Interface().(bool); ok && v {
					keepList = keepList[:len(keepList)-1]
					return
				}
				// log.Debugf("%v (%T)\n", script, script)
				
			// A function call that returns two values, the first a result value, and the second an error value
			// Requires one argument, no type restrictions, so long as the client reads it properly.
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case singleArgDualReturnCall:
				ret, callErr := (script.fn.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(arg))
				if !callErr.IsNil() {
					log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
					return
				}
				if v, ok := ret.Interface().(bool); ok && v {
					keepList = keepList[:len(keepList)-1]
					return
				}
			default:
				log.Debugf("Unhandled task found: Type:%T, Value:%v\n", script.fn, script.fn)
				return
			}
		}()
	}
	s.RUnlock()
	wait.Wait()
	s.Lock()
	defer s.Unlock()
	// for _, script := range removeList {
		// if script.id >= len(s.ScriptCalls) {
			// return
		// }
		// log.Debugf("Removing %v (type %T) from task list!\n", script.id, script)
		// if script.id+1 < len(s.ScriptCalls) {
			// s.ScriptCalls = append(s.ScriptCalls[:script.id], s.ScriptCalls[script.id+1:]...)
		// } else {
			// s.ScriptCalls = s.ScriptCalls[:i]
		// }
	// }
	// for i, v := range keepList {
		// log.Debug(i, "at", v)
	// }
	// s.ScriptCalls = keepList
	// s.ScriptCalls = make(ScriptCalls, len(keepList))
	s.ScriptCalls = keepList
	// copy(s.ScriptCalls, keepList)
	// for _, script := range removeList {
		// for i := 0; i < len(s.ScriptCalls); i++ {
			// v := s.ScriptCalls[i]
			// if v == nil {
				// continue
			// }
			// if v == script {
				// s.ScriptCalls[i] = nil
				// for ; i < len(s.ScriptCalls); i++ {
					// if s.ScriptCalls[i] != nil {
						// s.ScriptCalls[i], s.ScriptCalls[i+1] = s.ScriptCalls[i+1], s.ScriptCalls[i]
					// }
				// }
				// break
			// }
		// }
	// }
}
/*
func (s *Scripts) Call(v interface{}) {
	wait := sync.WaitGroup{}
	removeList := make([]int, 0, len(s.ScriptCalls))
	s.RLock()
	for i, script := range s.ScriptCalls {
		wait.Add(1)
		go func(script scriptCall) {
			defer wait.Done()
			// Determine the type of our script callback
			switch script.(type) {
			// Simple function call, no input no input
			case call:
				(script.(call))()
				// keep = append(keep, script)
			// A function call taking a *world.Player as an argument
			case playerArgCall:
				(script.(playerArgCall))(v.(entity.MobileEntity))
				// keep = append(keep, script)
			// A function call returning its active status.
			case StatusReturnCall:
				if (script.(StatusReturnCall))() {
					removeList = append(removeList, i)
					// keep = append(keep, script)
				}
			// A function call taking a *world.Player as an argument and returning its active status.
			case playerArgStatusReturnCall:
				if (script.(playerArgStatusReturnCall))(v.(entity.MobileEntity)) {
					removeList = append(removeList, i)
					// keep = append(keep, script)
					// removing = append(removing, i)
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case dualReturnCall:
				ret, callErr := (script.(dualReturnCall))(context.Background())
				if !callErr.IsNil() {
					removeList = append(removeList, i)
					log.Warn("Error retVal from a dualReturnCall in the Anko ctx:", callErr.Elem())
					return
				}
				if v, ok := ret.Interface().(bool); ok && v {
					removeList = append(removeList, i)
					// removing = append(removing, i)
					// keep = append(keep, script)
					return
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Requires one argument, no type restrictions, so long as the client reads it properly.
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case singleArgDualReturnCall:
				ret, callErr := (script.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(v))
				if !callErr.IsNil() {
					removeList = append(removeList, i)
					// removing = append(removing, i)
					log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
					return
				}
				if v, ok := ret.Interface().(bool); ok && v {
					removeList = append(removeList, i)
					// removing = append(removing, i)
					// keep = append(keep, script)
					return
				}
			default:
				return
			}
		}(script)
	}
	s.RUnlock()
	wait.Wait()
	s.Lock()
	defer s.Unlock()
	for _, i := range removeList {
		size := len(s.ScriptCalls)
		if i < size-1 {
			s.ScriptCalls = append(s.ScriptCalls[:i], s.ScriptCalls[i+1:]...)
		} else {
			s.ScriptCalls = s.ScriptCalls[:i]
		}
	}
	// s.ScriptCalls = make(ScriptCalls, len(keep))
	// copy(s.ScriptCalls, keep)
	// s.ScriptCalls = keep
}
*/
