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
	"errors"
	"sync"
	"sync/atomic"
)

const (
	chanOpened = iota
	chanClosed
)

// Chan is a concurrent channel to avoid the panic when sending data to a closed channel.
// It is safe for multiple goroutines to call a Chan's methods concurrently.
type Chan struct {
	// status is used to mark the status of channel.
	// It is set to chanOpened after created, and set to chanClosed before close.
	status int32

	// ch is the real channel of Chan.
	ch chan interface{}

	// wg is used to count sending data.
	// Before data is sent, do wg.Add(1).
	// After data is sent, do wg.Done().
	// Before close channel, do wg.Wait() to make sure all sending data are sent to channel.
	wg sync.WaitGroup
}

// Create a new SafeChan. If rt is nil, means any type data can be send to it.
func MakeChan(size int) *Chan {
	if size < 0 {
		return nil
	}
	sc := &Chan{
		status: chanOpened,
		ch:     make(chan interface{}, size),
	}
	return sc
}

// Get the read-only channel
func (c *Chan) Chan() <-chan interface{} {
	return c.ch
}

// Check if the channel is closed
func (c *Chan) IsClosed() bool {
	return atomic.LoadInt32(&c.status) != chanOpened
}

// Push data into channel
func (c *Chan) Push(data interface{}) error {
	if c.IsClosed() {
		return errors.New("channel is closed")
	}
	c.wg.Add(1)
	c.ch <- data
	c.wg.Done()
	return nil
}

// Get data from channel
func (c *Chan) Get() interface{} {
	return <-c.ch
}

// Close channel
func (c *Chan) Close() {
	// mark channel as closed
	if !atomic.CompareAndSwapInt32(&c.status, chanOpened, chanClosed) {
		return
	}
	go closeChan(c)
}

func closeChan(c *Chan) {
	// avoid panic when send data
	c.wg.Wait()
	close(c.ch)
}
