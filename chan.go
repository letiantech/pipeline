package pipeline

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

type SafeChan interface {
	Chan() <-chan interface{}
	Type() reflect.Type
	Push(data interface{}) error
	Get() interface{}
	IsClosed() bool
	Close()
}

//SafeChan is used to avoid the panic when sending data to a closed channel
type safeChan struct {
	closed int32
	ch     chan interface{}
	rt     reflect.Type
	wg     sync.WaitGroup
}

//Create a new SafeChan. If rt is nil, means any type data can be send to it
func NewSafeChan(size int, rt reflect.Type) SafeChan {
	if size < 0 {
		return nil
	}
	sc := &safeChan{
		closed: 0,
		ch:     make(chan interface{}, size),
		rt:     rt,
	}
	return sc
}

//Get the real channel
func (c *safeChan) Chan() <-chan interface{} {
	return c.ch
}

//Get the data type of channel
func (c *safeChan) Type() reflect.Type {
	return c.rt
}

//Check if the channel is closed
func (c *safeChan) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) != 0
}

//Push data into channel
func (c *safeChan) Push(data interface{}) error {
	if c.IsClosed() {
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
func (c *safeChan) Get() interface{} {
	return <-c.ch
}

//Close channel
func (c *safeChan) Close() {
	//mark channel as closed
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}
	go closeChan(c)
}

func closeChan(c *safeChan) {
	// avoid panic when send data
	c.wg.Wait()
	close(c.ch)
}
