package pipeline

import "sync/atomic"

type Task struct {
	finished int32
	sync     bool
	data     atomic.Value
	ch       chan interface{}
}

func NewTask(data interface{}, sync bool) *Task {
	t := &Task{
		sync: sync,
	}
	t.data.Store(data)
	if t.sync {
		t.ch = make(chan interface{}, 1)
	}
	return t
}

func (t *Task) Data() interface{} {
	return t.data.Load()
}

func (t *Task) Update(data interface{}) {
	t.data.Store(data)
}

func (t *Task) Finish() {
	atomic.StoreInt32(&t.finished, 1)
	if t.sync {
		t.ch <- t.data.Load()
	}
}

func (t *Task) Wait() interface{} {
	if !t.sync || atomic.LoadInt32(&t.finished) <= 0 {
		return nil
	}
	return <-t.ch
}
