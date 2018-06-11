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

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Data interface {
	Tag() string
}

type Task interface {
	Update(Data)
	GetData() Data
	Finish()
	Finished() <-chan struct{}
	Cancel()
	Canceled() <-chan struct{}
	Tag() string
}

type task struct {
	sync.RWMutex
	Data
	ctx        context.Context
	cancel     context.CancelFunc
	isFinished int32
	isCanceled int32
	finished   chan struct{}
}

// Create new task by data.
//
// Example:
//  //create new task
//  t := NewTask(data, 100 * time.Millisecond)
func NewTask(data Data, timeout time.Duration) Task {
	if data == nil {
		return nil
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	t := &task{
		ctx:        ctx,
		cancel:     cancelFunc,
		Data:       data,
		finished:   make(chan struct{}, 0),
		isFinished: 0,
		isCanceled: 0,
	}
	return t
}

//get the data of task
func (t *task) GetData() Data {
	t.RLock()
	defer t.RUnlock()
	return t.Data
}

//update the data of task
func (t *task) Update(data Data) {
	t.Lock()
	defer t.Unlock()
	t.Data = data
}

func (t *task) Cancel() {
	if !atomic.CompareAndSwapInt32(&t.isCanceled, 0, 1) {
		return
	}
	t.cancel()
}

func (t *task) Canceled() <-chan struct{} {
	if _, ok := t.ctx.Deadline(); ok {
		atomic.CompareAndSwapInt32(&t.isCanceled, 0, 1)
	}
	return t.ctx.Done()
}

func (t *task) Finish() {
	if !atomic.CompareAndSwapInt32(&t.isFinished, 0, 1) {
		return
	}
	close(t.finished)
}

func (t *task) Finished() <-chan struct{} {
	return t.finished
}
