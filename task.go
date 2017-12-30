package pipeline

import "sync/atomic"

//Task is used to perform a sync request in pipeline
type Task struct {
	finished int32
	ch       SafeChan
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
		ch:       NewSafeChan(1, nil),
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
