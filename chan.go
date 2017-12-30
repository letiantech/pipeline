package pipeline

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

//SafeChan is used to avoid the panic when sending data to a closed channel
type SafeChan struct {
	closed int32
	ch     chan interface{}
	rt     reflect.Type
	wg     sync.WaitGroup
}

//Create a new SafeChan. If rt is nil, means any type data can be send to it
func NewSafeChan(size int, rt reflect.Type) *SafeChan {
	if size < 0 {
		return nil
	}
	sc := &SafeChan{
		closed: 0,
		ch:     make(chan interface{}, size),
		rt:     rt,
	}
	return sc
}

//Get the real channel
func (c *SafeChan) Chan() <-chan interface{} {
	return c.ch
}

//Get the data type of channel
func (c *SafeChan) Type() reflect.Type {
	return c.rt
}

//Check if the channel is closed
func (c *SafeChan) Closed() bool {
	return atomic.LoadInt32(&c.closed) != 0
}

//Push data into channel
func (c *SafeChan) Push(data interface{}) error {
	if c.Closed() {
		return errors.New("channel is closed")
	}
	if c.rt != nil && reflect.TypeOf(data) != c.rt {
		return errors.New("data type different with channel type")
	}
	c.wg.Add(1)
	c.ch <- data
	c.wg.Done()
	return nil
}

//Get data from channel
func (c *SafeChan) Get() interface{} {
	return <-c.ch
}

//Close channel
func (c *SafeChan) Close() {
	if c.Closed() {
		return
	}
	//mark channel as closed
	atomic.StoreInt32(&c.closed, 1)
	go closeChan(c)
}

func closeChan(c *SafeChan) {
	// avoid panic when send data
	c.wg.Wait()
	close(c.ch)
}
