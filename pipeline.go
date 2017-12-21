package gpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Closed  = 0
	Running = 1
	Stopped = 2
)

type pipeLine struct {
	wg   sync.WaitGroup
	stat int32
	buf  chan interface{}
	size int
}

type PipeLine interface {
	Start()
	In(data interface{})
	Out() interface{}
	Stop()
	Close()
}

func NewPipeLine(size int) (PipeLine, error) {
	if size < 1 {
		return nil, errors.New("pipe line size must greater than 0")
	}
	p := &pipeLine{
		buf:  make(chan interface{}, size),
		size: size,
	}
	p.Start()
	return PipeLine(p), nil
}

func (p *pipeLine) Start() {
	atomic.StoreInt32(&p.stat, Running)
}

func (p *pipeLine) Stop() {
	atomic.StoreInt32(&p.stat, Stopped)
}

func (p *pipeLine) Close() {
	p.Stop()
	go func() {
		defer p.wg.Done()
		p.wg.Add(1)
		CloseChan(p.buf)
	}()
	p.wg.Wait()
	atomic.StoreInt32(&p.stat, Closed)
}

func (p *pipeLine) In(data interface{}) {
	if atomic.LoadInt32(&p.stat) != Running {
		return
	}
	p.buf <- data
}

func (p *pipeLine) Out() interface{} {
	if atomic.LoadInt32(&p.stat) != Running {
		return nil
	}
	return <-p.buf
}

func CloseChan(ch chan interface{}) {
	duration := 50 * time.Millisecond
	t := time.NewTimer(duration)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
			t.Reset(duration)
			break
		case <-t.C:
			close(ch)
			return
		}
	}
}
