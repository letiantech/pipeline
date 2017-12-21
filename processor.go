package gpool

import (
	"sync"
	"sync/atomic"
	"time"
)

type Processor struct {
	wg    sync.WaitGroup
	stat  int32
	in    PipeLine
	out   PipeLine
	speed int32
	f     func(data interface{}) interface{}
}

type ProcessorConfig struct {
	Speed float32
	In    PipeLine
	Out   PipeLine
	Func  func(data interface{}) interface{}
}

func NewProcessor(cfg *ProcessorConfig) *Processor {
	if cfg.In == nil || cfg.Out == nil {
		return nil
	}
	l := &Processor{
		in:    cfg.In,
		out:   cfg.Out,
		speed: int32(cfg.Speed * 1000),
		f:     cfg.Func,
	}
	l.Start()
	return l
}

func (l *Processor) SetSpeed(speed float64) {
	atomic.StoreInt32(&l.speed, int32(speed*1000))
}

func (l *Processor) Start() {
	atomic.StoreInt32(&l.stat, Running)
	go run(l)
}

func (l *Processor) Stop() {
	atomic.StoreInt32(&l.stat, Stopped)
	l.wg.Wait()
}

func (l *Processor) Close() {
	l.Stop()
	atomic.StoreInt32(&l.stat, Closed)
}

func (l *Processor) CloseAll() {
	l.Close()
	l.in.Close()
	l.out.Close()
	atomic.StoreInt32(&l.stat, Closed)
}

func (l *Processor) In(data interface{}) {
	if atomic.LoadInt32(&l.stat) != Running || l.in == nil {
		return
	}
	l.in.In(data)
}

func (l *Processor) Out() interface{} {
	if atomic.LoadInt32(&l.stat) != Running || l.out == nil {
		return nil
	}
	return l.out.Out()
}

func run(l *Processor) {
	defer l.wg.Done()
	l.wg.Add(1)
	var duration = time.Duration(0)
	var data interface{}
	for {
		speed := atomic.LoadInt32(&l.speed)
		if speed > 0 {
			duration = (time.Second / time.Duration(speed)) * 1000
			data = l.in.Out()
			if data == nil || atomic.LoadInt32(&l.stat) != Running {
				return
			}
			if l.f != nil {
				data = l.f(data)
			}
			if l.out != nil {
				l.out.In(data)
			}
			if duration > 0 {
				time.Sleep(duration)
			}
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}
