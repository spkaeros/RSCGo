 
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
	scriptCall interface{}
	//scriptArg This is a type-free box type, made for boxing any function arguments, for use in conjunction with scriptCall.
	scriptArg interface{}
	//scriptArg A slice of scriptCalls.
	scriptCalls []scriptCall
	//Scripts A locked slice of scriptCalls.
	Scripts struct {
		scriptCalls
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
		curTick++
		if curTick >= ticks {
			curTick = 0
			return fn()
		}
		return false
	}
	s.Add(ticker)
}

func (s *Scripts) Add(fn scriptCall) {
	s.Lock()
	defer s.Unlock()
	s.scriptCalls = append(s.scriptCalls, fn)
}

func (s *Scripts) Tick() {
	s.Call(nil)
}

func (s *Scripts) Call(v interface{}) {
	wait := sync.WaitGroup{}
	keep := make(scriptCalls, 0, len(s.scriptCalls))
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
				keep = append(keep, script)
			// A function call taking a *world.Player as an argument
			case playerArgCall:
				(script.(playerArgCall))(v.(entity.MobileEntity))
				keep = append(keep, script)
			// A function call returning its active status.
			case StatusReturnCall:
				if !(script.(StatusReturnCall))() {
					// removing = append(removing, i)
					keep = append(keep, script)
				}
			// A function call taking a *world.Player as an argument and returning its active status.
			case playerArgStatusReturnCall:
				if !(script.(playerArgStatusReturnCall))(v.(entity.MobileEntity)) {
					keep = append(keep, script)
					// removing = append(removing, i)
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
					// removing = append(removing, i)
					keep = append(keep, script)
					return
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Requires one argument, no type restrictions, so long as the client reads it properly.
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case singleArgDualReturnCall:
				ret, callErr := (script.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(v))
				if !callErr.IsNil() {
					// removing = append(removing, i)
					log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
					return
				}
				if v, ok := ret.Interface().(bool); ok && !v {
					// removing = append(removing, i)
					keep = append(keep, script)
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
	copy(s.scriptCalls, keep)
	
	// s.scriptCalls = keep
}
