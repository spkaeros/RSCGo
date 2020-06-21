 
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
	
	`github.com/spkaeros/rscgo/pkg/config`
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

func Schedule(ticks int, call func() bool) {
	ticksLeft := ticks
	TickList.Add(func() bool {
		ticksLeft -= 1
		if ticksLeft <= 0 {
			ticksLeft = ticks
			if call() {
				return true
			}
		}
		return false
	})
}

func (s *Scripts) Schedule(ticks int, fn func() bool) {
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

//Range runs fn(Task) for each Task in the list.
func (t *TaskList) Range(fn func(Task)) {
	t.RLock()
	defer t.RUnlock()
	for _, task := range t.Tasks {
		if fn == nil {
			log.Debug("niltask found!")
			continue
		}
		fn(task)
	}
}

//RunSynchronous will execute every task in the list sequentially, one at a time.  The order in which they run is
// basically unpredictable.  Upon completion, the tasks that have returned true will be removed from the set, then
// the calling function will resume execution.
func (t *TaskList) RunSynchronous() Tasks {
	t.Lock()
	defer t.Unlock()
	retainedTasks := []Task{}
	for _, task := range t.Tasks {
		if task() {
			if config.Verbosity >= 2 {
				log.Debugf("task finished; removing from collection.\n")
			}
		} else {
			retainedTasks = append(retainedTasks, task)
		}
	}
//	t.Tasks = append(t.Tasks[:0], retainedTasks[:]...)
	return retainedTasks
}

//RunAsynchronous will attempt to run every task in the list at the same time, but not necessarily in parallel
// (it depends on the hardware available, more cores means more parallelism), and then waits for them to complete.
// No guarantees are offered on what order the tasks will execute, only that they will all be executed.
// Upon completion of all tasks, it will remove any tasks that had returned true, then the caller will resume execution.
// TODO: Probably some form of pooling?
func (t *TaskList) RunAsynchronous() {
	wait := sync.WaitGroup{}
	// retainedTasks := make(chan Task, len(t.Tasks))
	t.RLock()
	for _, task := range t.Tasks {
		wait.Add(1)
		go func(task Task) {
			defer wait.Done()
			if task == nil {
				return
			}
			if task() {
				log.Debugf("task finished; removing from collection.\n")
			} else {
				// retainedTasks <- task
				t.Lock()
				t.Tasks = append(t.Tasks, task)
				t.Unlock()
			}
		}(task)
	}
	t.RUnlock()
	wait.Wait()
	t.Lock()
	defer t.Unlock()
	t.Tasks = t.Tasks[:0]
	// select {
	// case task, ok := <-retainedTasks:
		// if !ok {
			// break
		// }
		// 
	// default: return
	// }
//	t.tasks = append(t.tasks[:0], retainedTasks[:]...)
	return
}

//Count returns the total number of tasks currently in this list.
func (t *TaskList) Count() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.Tasks)
}

//Add will add a mapping of the Task fn to the provided name, in this list.
func (t *TaskList) Add(fn Task) int {
	t.Lock()
	defer t.Unlock()
	t.Tasks = append(t.Tasks, fn)
	return len(t.Tasks)-1
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
	removing := make([]int, 0, len(s.scriptCalls))
	s.RLock()
	for i, script := range s.scriptCalls {
		wait.Add(1)
		go func(script scriptCall) {
			defer wait.Done()
			// Determine the type of our script callback
			switch script.(type) {
			// Simple function call, no input no input
			case call:
				(script.(call))()
			// A function call taking a *world.Player as an argument
			case playerArgCall:
				(script.(playerArgCall))(v.(entity.MobileEntity))
			// A function call returning its active status.
			case StatusReturnCall:
				if (script.(StatusReturnCall))() {
					removing = append(removing, i)
				}
			// A function call taking a *world.Player as an argument and returning its active status.
			case playerArgStatusReturnCall:
				if (script.(playerArgStatusReturnCall))(v.(entity.MobileEntity)) {
					removing = append(removing, i)
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
				if v, ok := ret.Interface().(bool); ok && v {
					log.Debugf("%v, %v\n", v, callErr)
					removing = append(removing, i)
					return
				}
			// A function call that returns two values, the first a result value, and the second an error value
			// Requires one argument, no type restrictions, so long as the client reads it properly.
			// Upon non-nil error value, it will log the stringified err struct then remove from active list,
			// otherwise schedules the same call to run again next tick.
			case singleArgDualReturnCall:
				ret, callErr := (script.(singleArgDualReturnCall))(context.Background(), reflect.ValueOf(v))
				if !callErr.IsNil() {
					removing = append(removing, i)
					log.Warn("Error retVal from a singleArgDualReturnCall in the Anko ctx:", callErr.String())
					return
				}
				if v, ok := ret.Interface().(bool); ok && v {
					removing = append(removing, i)
					log.Debug(ret)
					return
				}
			default: return
			}
		}(script)
	}
	s.RUnlock()
	wait.Wait()
	s.Lock()
	defer s.Unlock()
	for _, v := range removing {
		s.scriptCalls = s.scriptCalls[:v]
		if v < len(s.scriptCalls)-1 {
			s.scriptCalls = append(s.scriptCalls[:v], s.scriptCalls[v+1:])
		}
	}
}

//Get returns the Task at the given index of this collection.
// If no task is at the provided index, returns nil.
func (t *TaskList) Get(i int) Task {
	t.RLock()
	defer t.RUnlock()
	if i >= len(t.Tasks) {
		return nil
	}
	return t.Tasks[i]
}

//Remove locates and removes the given task, and returns its old index.
// If it finds no such task, returns -1
func (t *TaskList) Remove(fn Task) int {
	t.Lock()
	defer t.Unlock()
	for i, v := range t.Tasks {
		if &v == &fn {
			t.Tasks[i] = nil
			if i >= len(t.Tasks)-1 {
				t.Tasks = t.Tasks[:i]
				return i
			}
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			return i
		}
	}
	return -1
}

//Remove locates and removes the given task, and returns its old index.
// If it finds no such task, returns -1
func (t *TaskList) RemoveIdx(i int) int {
	t.Lock()
	defer t.Unlock()
	t.Tasks[i] = nil
	return i
}