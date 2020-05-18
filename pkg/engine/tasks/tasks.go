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

type tasks []Task

//TaskList is a concurrency-safe list of Task closures mapped to simple string identifiers.
// The main purpose for this is to provide the ability to schedule functions to run on the
// game engines thread, to synchronize certain sensitive events.
// It is guarded by a sync.RWMutex.
type TaskList struct {
	tasks
	sync.RWMutex
}

//TickList A collection of Tasks that are intended to be ran once per game engine tick.
// Tasks should contractually return either true if they are to be removed after execution completes,
// or false if they are to be ran again on the next engine cycle.
var TickList = &TaskList{}

//Range runs fn(Task) for each Task in the list.
func (t *TaskList) Range(fn func(Task)) {
	t.RLock()
	defer t.RUnlock()
	for _, task := range t.tasks {
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
func (t *TaskList) RunSynchronous() {
	t.Lock()
	defer t.Unlock()
	retainedTasks := []Task{}
	for _, task := range t.tasks {
		if task() {
			if config.Verbosity >= 2 {
				log.Debugf("task finished; removing from collection.\n")
			}
		} else {
			retainedTasks = append(retainedTasks, task)
		}
	}
	t.tasks = append(t.tasks[:0], retainedTasks[:]...)
}

//RunAsynchronous will attempt to run every task in the list at the same time, but not necessarily in parallel
// (it depends on the hardware available, more cores means more parallelism), and then waits for them to complete.
// No guarantees are offered on what order the tasks will execute, only that they will all be executed.
// Upon completion of all tasks, it will remove any tasks that had returned true, then the caller will resume execution.
// TODO: Probably some form of pooling?
func (t *TaskList) RunAsynchronous() {
	retainedTasks := []Task{}
	wait := sync.WaitGroup{}
	t.Range(func(task Task) {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if task() {
				if config.Verbosity >= 2 {
					log.Debugf("task finished; removing from collection.\n")
				}
			} else {
				retainedTasks = append(retainedTasks, task)
			}
		}()
	})
	wait.Wait()
	t.Lock()
	t.tasks = append(t.tasks[:0], retainedTasks[:]...)
	t.Unlock()
}

//Count returns the total number of tasks currently in this list.
func (t *TaskList) Count() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.tasks)
}

//Add will add a mapping of the Task fn to the provided name, in this list.
func (t *TaskList) Add(fn Task) int {
	t.Lock()
	defer t.Unlock()
	t.tasks = append(t.tasks, fn)
	return len(t.tasks)-1
}

//Get returns the Task at the given index of this collection.
// If no task is at the provided index, returns nil.
func (t *TaskList) Get(i int) Task {
	t.RLock()
	defer t.RUnlock()
	if i >= len(t.tasks) {
		return nil
	}
	return t.tasks[i]
}

//Remove locates and removes the given task, and returns its old index.
// If it finds no such task, returns -1
func (t *TaskList) Remove(fn Task) int {
	t.Lock()
	defer t.Unlock()
	for i, v := range t.tasks {
		if &v == &fn {
			t.tasks[i] = nil
			if i >= len(t.tasks)-1 {
				t.tasks = t.tasks[:i]
				return i
			}
			t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
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
	t.tasks[i] = nil
	return i
}
