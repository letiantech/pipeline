package pipeline

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const (
	stateOpened = iota
	stateClosed
)

const (
	stageSleeping = iota
	stageGetting
	stageExecuting
	stagePushing
	stageExiting
)

var _ Pipeline = &BasePipeline{}

var allCreators = make(map[string]func(*Config) (Pipeline, error))
var typeOfTask = reflect.TypeOf(Task{})

type Pipeline interface {
	Close()
	CloseAll()
	SetSpeed(speed float32)
	IsClosed() bool
	Type() reflect.Type
	Push(data interface{})
	Connect(next Pipeline)
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

type BasePipeline struct {
	cfg     Config
	next    sync.Map
	state   int32
	stage   int32
	buffer  *Chan
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
	p := &BasePipeline{}
	p.cfg = *cfg
	if p.cfg.BufferSize < 1 {
		return nil, errors.New("pipeline buffer size must greater than 0")
	}
	if p.cfg.PoolSize < 1 {
		return nil, errors.New("pipeline pool size must greater than 0")
	}
	p.buffer = MakeChan(cfg.BufferSize)
	p.next = sync.Map{}
	for _, c := range cfg.NextConfigs {
		np, err := NewPipeline(c)
		if err != nil {
			p.CloseAll()
			return nil, err
		}
		p.next.Store(c.Type, np)
	}
	p.state = stateOpened
	for i := 0; i < p.cfg.PoolSize; i++ {
		go p.run()
	}
	return p, nil
}

func (p *BasePipeline) Close() {
	if !atomic.CompareAndSwapInt32(&p.state, stateOpened, stateClosed) {
		return
	}
	p.buffer.Close()
	p.wg.Wait()
}

func (p *BasePipeline) CloseAll() {
	if p.IsClosed() {
		return
	}
	p.Close()
	p.next.Range(func(key, value interface{}) bool {
		switch v := value.(type) {
		case Pipeline:
			v.CloseAll()
		}
		return true
	})
}

func (p *BasePipeline) Connect(next Pipeline) {
}

func (p *BasePipeline) Type() reflect.Type {
	return p.cfg.Type
}

func (p *BasePipeline) Push(data interface{}) {
	if p.IsClosed() || p.dataType(data) != p.cfg.Type {
		return
	}
	p.buffer.Push(data)
}

func (p *BasePipeline) SetSpeed(speed float32) {
	p.limiter.SetSpeed(speed)
}

func (p *BasePipeline) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == stateClosed
}

func (p *BasePipeline) dataType(data interface{}) reflect.Type {
	rt := reflect.TypeOf(data)
	if rt == typeOfTask {
		t := data.(Task)
		rt = reflect.TypeOf(t.getData())
	}
	return rt
}

func (p *BasePipeline) run() {
	defer p.wg.Done()
	var data interface{}
	p.wg.Add(1)
	var stage = stageSleeping
	for {
		switch stage {
		case stageSleeping:
			if duration := p.limiter.Update(); duration > 0 {
				time.Sleep(duration)
			}
			stage = stageGetting
		case stageGetting:
			data = p.buffer.Get()
			stage = stageExecuting
			if data == nil {
				stage = stageExiting
			}
		case stageExecuting:
			if p.cfg.Func != nil {
				if reflect.TypeOf(data) == typeOfTask {
					t := data.(Task)
					t.update(p.cfg.Func(t.getData()))
					data = t
				} else {
					data = p.cfg.Func(data)
				}
			}
			stage = stagePushing
		case stagePushing:
			if reflect.TypeOf(data) == typeOfTask {
				t := data.(Task)
				var np Pipeline = nil
				val, ok := p.next.Load(reflect.TypeOf(t.getData()))
				if ok {
					switch v := val.(type) {
					case Pipeline:
						np = v
					}
				}
				if np != nil {
					np.Push(data)
				} else {
					t.finish()
				}
			} else {
				var np Pipeline = nil
				val, ok := p.next.Load(reflect.TypeOf(data))
				if ok {
					switch v := val.(type) {
					case Pipeline:
						np = v
					}
				}
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
		if p.IsClosed() || stage == stageExiting {
			return
		}
	}
}
