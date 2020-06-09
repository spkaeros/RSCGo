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
	"sync"
	
	`github.com/spkaeros/rscgo/pkg/config`
	`github.com/spkaeros/rscgo/pkg/log`
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
var TickList = &TaskList{}

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
	retainedTasks := make(chan Task, 255)
	t.Lock()
	for _, task := range t.Tasks {
		wait.Add(1)
		go func(task Task) {
			defer wait.Done()
			if task == nil {
				return
			}
			if task() {
				defer t.Remove(task)
			} else {
				retainedTasks <- task
			}
		}(task)
	}
	t.Unlock()
	wait.Wait()
	// t.Tasks = t.Tasks[:0]
	// select {
	// case task, ok := <-retainedTasks:
		// if !ok {
			// break
		// }
		// t.Tasks = append(t.Tasks, task)
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
