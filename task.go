// Copyright (c) 2017 letian0805@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pipeline

import "sync/atomic"

//Task is used to perform a sync request in pipeline
type Task struct {
	finished int32
	ch       *Chan
	data     interface{}
	auto     bool
}

// Create new task by data.
//
// Notice:
// If auto is true, the task will close the channel of task automatic when finished.
// If auto is false, you must close the channel by using t.Destroy()
//
// Example:
//  //create new task
//  t := NewTask(data, true)
//  //push the task into a pipeline
//  pl.Push(t)
//  //wait the result
//  result := <-t.Chan()  //or using  result := t.Wait()
func NewTask(data interface{}, auto bool) Task {
	t := Task{
		ch:       MakeChan(1),
		finished: 0,
		data:     data,
		auto:     auto,
	}
	return t
}

// Get the result channel of a task.
//
// Notice:
// The result channel of a task will be closed automatic when finished,
// so be careful when waiting multi channel by using select in a loop
func (t *Task) Chan() <-chan interface{} {
	return t.ch.Chan()
}

//Wait the result of a task
func (t *Task) Wait() interface{} {
	return t.ch.Get()
}

//Destroy a finished task
func (t *Task) Destroy() {
	t.ch.Close()
}

//get the data of task
func (t *Task) getData() interface{} {
	return t.data
}

//update the data of task
func (t *Task) update(data interface{}) {
	t.data = data
}

//finish task. If auto is true, close channel automatic.
func (t *Task) finish() {
	if !atomic.CompareAndSwapInt32(&t.finished, 0, 1) {
		return
	}
	t.ch.Push(t.data)
	if t.auto {
		t.ch.Close()
	}
}
