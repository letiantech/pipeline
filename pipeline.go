package pipeline

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const (
	StateClosed    = 0
	StateOpened    = 1
	StateDestroyed = 3
)

const (
	stageSleeping = iota
	stageGetting
	stageExecuting
	stagePushing
	stageExiting
)

var _ Pipeline = &pipeline{}

var allCreators = make(map[string]func(*Config) (Pipeline, error))
var typeOfTask = reflect.TypeOf(Task{})

type Task struct {
	finished bool
	ch       chan interface{}
	data     interface{}
}

type Pipeline interface {
	Open()
	Close()
	Destroy()
	DestroyAll()
	SetSpeed(speed float32)
	IsOpened() bool
	IsClosed() bool
	IsDestroyed() bool
	Type() reflect.Type
	Push(data interface{})
	PushTask(data interface{}) *Task
}

type Config struct {
	Name        string
	Speed       float32
	BufferSize  int
	PoolSize    int
	Type        reflect.Type
	Func        func(data interface{}) interface{}
	NextConfigs []*Config
}

type pipeline struct {
	cfg     Config
	next    map[reflect.Type]Pipeline
	stat    int32
	stage   int32
	buffer  chan interface{}
	wg      sync.WaitGroup
	limiter Limiter
}

func init() {
	Register("base", basePipelineCreator)
}

func Register(name string, creator func(config *Config) (Pipeline, error)) {
	allCreators[name] = creator
}

func NewPipeline(cfg *Config) (Pipeline, error) {
	creator, ok := allCreators[cfg.Name]
	if !ok {
		return nil, errors.New("no creator for pipeline " + cfg.Name)
	}
	return creator(cfg)
}

func basePipelineCreator(cfg *Config) (Pipeline, error) {
	p := &pipeline{}
	p.cfg = *cfg
	if p.cfg.BufferSize < 1 {
		return nil, errors.New("pipeline buffer size must greater than 0")
	}
	if p.cfg.PoolSize < 1 {
		return nil, errors.New("pipeline pool size must greater than 0")
	}
	p.buffer = make(chan interface{}, p.cfg.BufferSize)
	p.next = make(map[reflect.Type]Pipeline)
	for _, c := range cfg.NextConfigs {
		np, err := NewPipeline(c)
		if err != nil {
			p.DestroyAll()
			return nil, err
		}
		p.next[c.Type] = np
	}
	p.Open()
	for i := 0; i < p.cfg.PoolSize; i++ {
		go p.run()
	}
	return p, nil
}

func (p *pipeline) Open() {
	if p.IsOpened() {
		return
	}
	atomic.StoreInt32(&p.stat, StateOpened)
}

func (p *pipeline) Close() {
	if !p.IsOpened() {
		return
	}

	atomic.StoreInt32(&p.stat, StateClosed)
}

func (p *pipeline) Destroy() {
	if p.IsDestroyed() {
		return
	}
	atomic.StoreInt32(&p.stat, StateDestroyed)
	CloseChan(p.buffer)
	p.wg.Wait()
}

func (p *pipeline) DestroyAll() {
	if p.IsDestroyed() {
		return
	}
	p.Destroy()
	for _, v := range p.next {
		v.DestroyAll()
	}
}

func (p *pipeline) Type() reflect.Type {
	return p.cfg.Type
}

func (p *pipeline) Push(data interface{}) {
	if p.IsDestroyed() || p.dataType(data) != p.cfg.Type {
		return
	}
	p.buffer <- data
}

func (p *pipeline) PushTask(data interface{}) *Task {
	if p.IsDestroyed() || reflect.TypeOf(data) != p.cfg.Type {
		return nil
	}
	t := NewTask(data)
	p.buffer <- t
	return &t
}

func (p *pipeline) SetSpeed(speed float32) {
	p.limiter.SetSpeed(speed)
}

func (p *pipeline) IsClosed() bool {
	return atomic.LoadInt32(&p.stat) == StateClosed
}

func (p *pipeline) IsOpened() bool {
	return atomic.LoadInt32(&p.stat) == StateOpened
}

func (p *pipeline) IsDestroyed() bool {
	return atomic.LoadInt32(&p.stat) == StateDestroyed
}

func (p *pipeline) dataType(data interface{}) reflect.Type {
	rt := reflect.TypeOf(data)
	if rt == typeOfTask {
		t := data.(Task)
		rt = reflect.TypeOf(t.GetData())
	}
	return rt
}

func (p *pipeline) run() {
	defer p.wg.Done()
	var data interface{}
	p.wg.Add(1)
	var stage = stageSleeping
	for {
		switch stage {
		case stageSleeping:
			duration := 10 * time.Millisecond
			if p.IsOpened() {
				duration = p.limiter.Update()
			}
			if duration > 0 {
				time.Sleep(duration)
			}
			stage = stageGetting
		case stageGetting:
			data = <-p.buffer
			stage = stageExecuting
			if data == nil {
				stage = stageExiting
			}
		case stageExecuting:
			if p.cfg.Func != nil {
				if reflect.TypeOf(data) == typeOfTask {
					t := data.(Task)
					t.update(p.cfg.Func(t.GetData()))
					data = t
				} else {
					data = p.cfg.Func(data)
				}
			}
			stage = stagePushing
		case stagePushing:
			if reflect.TypeOf(data) == typeOfTask {
				t := data.(Task)
				np := p.next[reflect.TypeOf(t.GetData())]
				if np != nil {
					np.Push(data)
				} else {
					t.finish()
				}
			} else {
				np := p.next[reflect.TypeOf(data)]
				if np != nil {
					np.Push(data)
				}
			}
			stage = stageSleeping
		case stageExiting:
			return
		default:
			stage = stageSleeping
		}
		if p.IsDestroyed() || stage == stageExiting {
			return
		}
		if p.IsClosed() {
			stage = stageSleeping
		}
	}
}

func CloseChan(ch chan interface{}) {
	duration := 10 * time.Millisecond
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
