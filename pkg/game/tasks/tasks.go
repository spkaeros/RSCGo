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
)

//Task is a single func that takes no args and returns a bool to indicate whether or not it
// should be removed from the set.
type Task func() bool

//taskSet is a non-concurrent collection type that maps Tasks to their identifying names.
// This is not intended for direct usage, but to be composed with concurrent locking primitives,
// and such, to create useful types for concurrent storing and retrieval of Task closures.
// See TaskList for an example of such a type.
type taskSet map[string]Task

//TaskList is a concurrency-safe list of Task closures mapped to simple string identifiers.
// The main purpose for this is to provide the ability to schedule functions to run on the
// game engines thread, to synchronize certain sensitive events.
// It is guarded by a sync.RWMutex.
type TaskList struct {
	taskSet
	sync.RWMutex
}

//TickerList A pointer to a TaskList collection to provide concurrent access to a list of
// Task closures that are intended to be ran once per game engine tick.
// Task closures return either true if they are to be removed immediately after they terminate
// from this collection, or false if they are to be ran again on the next engine cycle.
var TickerList = &TaskList{
	taskSet: make(taskSet),
}

//Range run fn(name,Task) for each Task in the receiving lists set.
func (t *TaskList) Range(fn func(string, Task)) {
	t.RLock()
	for name, task := range t.taskSet {
		fn(name, task)
	}
	t.RUnlock()
}

//RunSynchronous will run every task in the receiver task list one by one in a mostly
// unpredictable order.  Upon completion, the tasks that have returned true will be removed
// from the set and the calling function will resume its execution.
func (t *TaskList) RunSynchronous() {
	var removed []string
	t.Range(func(name string, task Task) {
		if task() {
			removed = append(removed, name)
		}
	})
	for _, taskName := range removed {
		t.Remove(taskName)
	}
}

//RunAsynchronous will attempt to run every task in the receiver task list concurrently, in an
// unpredictable order.  Upon completion, it waits for all tasks to finish execution,
// and then removes any tasks that had returned true.
// TODO: Probably some form of pooling?
func (t *TaskList) RunAsynchronous() {
	var removed []string
	var wg sync.WaitGroup
	wg.Add(t.Count())
	t.Range(func(name string, task Task) {
		//start := time.Now()
		defer wg.Done()
		go func() {

			if task() {
				removed = append(removed, name)
			}
		}()
		//log.Info.Printf("tickTask--%s; finished executing in %v", name, time.Since(start))
	})
	wg.Wait()
	for _, taskName := range removed {
		t.Remove(taskName)
	}
}

//Count returns the total number of tasks currently in the receiving lists set.
func (t *TaskList) Count() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.taskSet)
}

//Add adds a task to the receiving task list, mapped to the given string.
func (t *TaskList) Add(name string, fn Task) {
	t.Lock()
	defer t.Unlock()
	t.taskSet[name] = fn
}

//Get takes a task name as input, returns the task assigned to the given name as output.
// Returns null if there is no task with such a name in the receiving lists set.
func (t *TaskList) Get(name string) Task {
	t.RLock()
	defer t.RUnlock()
	return t.taskSet[name]
}

//Remove takes a task name as input and if any task exists with such a name in the
// receiving lists set, removes it from the set.
func (t *TaskList) Remove(name string) {
	t.Lock()
	defer t.Unlock()
	delete(t.taskSet, name)
}
