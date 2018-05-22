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
	"time"
)

type Sink interface {
	Push(data interface{}) error
	SetSpeed(speed float32)
}

type Source interface {
	Pull() interface{}
	SetSpeed(speed float32)
}

type Pipe interface {
	Init(size int) Pipe
	GetSource() Source
	GetSink() Sink
	Close()
}

type source struct {
	ch *Chan
	*Limiter
}

func (s *source) Pull() interface{} {
	if duration := s.Update(); duration > 0 {
		time.Sleep(duration)
	}
	return s.ch.Pull()
}

type sink struct {
	ch *Chan
	*Limiter
}

func (s *sink) Push(data interface{}) error {
	if duration := s.Update(); duration > 0 {
		time.Sleep(duration)
	}
	return s.ch.Push(data)
}

type BasePipe struct {
	ch  *Chan
	in  Limiter
	out Limiter
}

func (bp *BasePipe) Init(size int) Pipe {
	if bp.ch == nil {
		bp.ch = MakeChan(size)
	}
	return bp
}

func (bp *BasePipe) GetSource() Source {
	s := &source{}
	s.ch = bp.ch
	s.Limiter = &bp.out
	return s
}

func (bp *BasePipe) GetSink() Sink {
	s := &sink{}
	s.ch = bp.ch
	s.Limiter = &bp.in
	return s
}

func (bp *BasePipe) Close() {
	if bp.ch != nil {
		bp.ch.Close()
		bp.ch = nil
	}
}
