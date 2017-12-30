package pipeline

import "time"

func NewTask(data interface{}) Task {
	t := Task{
		ch:       make(chan interface{}, 1),
		finished: false,
		data:     data,
	}
	return t
}

func (t *Task) Chan() <-chan interface{} {
	return t.ch
}

func (t *Task) WaitData() interface{} {
	return <-t.ch
}

func (t *Task) Destroy() {
	close(t.ch)
}

func (t *Task) update(data interface{}) {
	t.data = data
}

func (t *Task) GetData() interface{} {
	return t.data
}

func (t *Task) finish() {
	t.finished = true
	t.ch <- t.data
	go func() {
		time.Sleep(200 * time.Millisecond)
		close(t.ch)
	}()
}
