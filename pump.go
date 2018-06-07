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
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	pumpStatusStopped = iota
	pumpStatusRunning
)

type Filter func(in Data) Data

func DefaultFilter(in Data) Data {
	return in
}

type Pump interface {
	Init(source Source) Pump
	Start()
	Stop()
	AddFilter(filter Filter, tag string)
	AddSink(sink Sink, tag string)
}

type BasePump struct {
	sync.RWMutex
	status  int32
	source  Source
	filters map[string]Filter
	sinks   map[string]Sink
}

func (bp *BasePump) Init(source Source) Pump {
	bp.source = source
	bp.sinks = make(map[string]Sink)
	bp.filters = make(map[string]Filter)
	return bp
}

func (bp *BasePump) Start() {
	if !atomic.CompareAndSwapInt32(&bp.status, pumpStatusStopped, pumpStatusRunning) {
		return
	}
	go bp.run()
}

func (bp *BasePump) Stop() {
	atomic.StoreInt32(&bp.status, pumpStatusStopped)
}

func (bp *BasePump) AddSink(sink Sink, tag string) {
	bp.Lock()
	defer bp.Unlock()
	bp.sinks[tag] = sink
}

func (bp *BasePump) AddFilter(filter Filter, tag string) {
	bp.Lock()
	defer bp.Unlock()
	bp.filters[tag] = filter
}

func (bp *BasePump) run() {
	defer bp.Stop()
	for {
		data := bp.source.Pull()
		if data == nil {
			break
		}
		bp.RLock()
		filter, ok := bp.filters[data.Tag()]
		bp.RUnlock()
		if !ok {
			filter = DefaultFilter
		}
		out := filter(data)
		if out != nil {
			bp.RLock()
			s, ok := bp.sinks[out.Tag()]
			bp.RUnlock()
			if !ok {
				continue
			}
			err := s.Push(out)
			if err != nil {
				fmt.Println(err)
			}
			var _ = err
		}
		if atomic.LoadInt32(&bp.status) != pumpStatusRunning {
			break
		}
	}
}
