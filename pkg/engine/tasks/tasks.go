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
// should be removed from the set it belongs to upon completion
type Task func() bool

//TaskList is a concurrency-safe list of Task closures mapped to simple string identifiers.
// The main purpose for this is to provide the ability to schedule functions to run on the
// game engines thread, to synchronize certain sensitive events.
// It is guarded by a sync.RWMutex.
type TaskList struct {
	taskSet map[string]Task
	sync.RWMutex
}

//Tickers A collection of Tasks that are intended to be ran once per game engine tick.
// Tasks should contractually return either true if they are to be removed after execution completes,
// or false if they are to be ran again on the next engine cycle.
var Tickers = &TaskList{
	taskSet: make(map[string]Task),
}

//Range runs fn(name,Task) for each Task in the list.
func (t *TaskList) Range(fn func(string, Task)) {
	t.RLock()
	for name, task := range t.taskSet {
		fn(name, task)
	}
	t.RUnlock()
}

//RunSynchronous will execute every task in the list sequentially, one at a time.  The order in which they run is
// basically unpredictable.  Upon completion, the tasks that have returned true will be removed from the set, then
// the calling function will resume execution.
func (t *TaskList) RunSynchronous() {
	var removeList []string
	t.Range(func(name string, task Task) {
		if task() {
			removeList = append(removeList, name)
		}
	})
	for _, taskName := range removeList {
		t.Remove(taskName)
	}
}

//RunAsynchronous will attempt to run every task in the list at the same time, but not necessarily in parallel
// (it depends on the hardware available, more cores means more parallelism), and then waits for them to complete.
// No guarantees are offered on what order the tasks will execute, only that they will all be executed.
// Upon completion of all tasks, it will remove any tasks that had returned true, then the caller will resume execution.
// TODO: Probably some form of pooling?
func (t *TaskList) RunAsynchronous() {
	var removeList []string
	var runningTasks sync.WaitGroup
	runningTasks.Add(t.Count())
	t.Range(func(name string, task Task) {
		// start := time.Now()
		defer runningTasks.Done()
		go func() {
			if task() {
				removeList = append(removeList, name)
			}
		}()
		// log.Info.Printf("tickTask--%s; finished executing in %v", name, time.Since(start))
	})
	runningTasks.Wait()
	for _, taskName := range removeList {
		t.Remove(taskName)
	}
}

//Count returns the total number of tasks currently in this list.
func (t *TaskList) Count() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.taskSet)
}

//Add will add a mapping of the Task fn to the provided name, in this list.
func (t *TaskList) Add(name string, fn Task) {
	t.Lock()
	defer t.Unlock()
	t.taskSet[name] = fn
}

//Get returns the Task that is mapped to the provided name in this list, or null if no such task exists.
func (t *TaskList) Get(name string) Task {
	t.RLock()
	defer t.RUnlock()
	return t.taskSet[name]
}

//Remove will remove the task mapped to the provided name from the list, if any such task exists.
func (t *TaskList) Remove(name string) {
	t.Lock()
	defer t.Unlock()
	delete(t.taskSet, name)
}
