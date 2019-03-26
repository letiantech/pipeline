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

package pipeline_test

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/letiantech/pipeline"
)

type tA struct {
	s string
}

func (a *tA) Tag() string {
	return reflect.TypeOf(a).String()
}

type tB struct {
	s string
}

func (b *tB) Tag() string {
	return reflect.TypeOf(b).String()
}

func dataFilter(data pipeline.Data) pipeline.Data {
	switch v := data.(type) {
	case *tB:
		return &tA{
			s: v.s,
		}
	case *tA:
		return &tB{
			s: v.s,
		}
	}
	return nil
}

func TestPump(t *testing.T) {
	const speedIn = 100
	const speedOut = 500
	const count = 300
	pipes := []pipeline.Pipe{
		new(pipeline.BasePipe).Init(10),
		new(pipeline.BasePipe).Init(10),
		new(pipeline.BasePipe).Init(10),
	}
	sources := make([]pipeline.Source, len(pipes))
	sinks := make([]pipeline.Sink, len(pipes))
	for i := 0; i < len(pipes); i++ {
		sources[i] = pipes[i].GetSource()
		sinks[i] = pipes[i].GetSink()
		sources[i].SetSpeed(speedOut)
		sinks[i].SetSpeed(speedIn)
	}
	at := (&tA{}).Tag()
	bt := (&tB{}).Tag()
	p := new(pipeline.BasePump).Init(sources[0])
	p.AddSink(sinks[1], at)
	p.AddSink(sinks[2], bt)
	p.AddFilter(dataFilter, at)
	p.AddFilter(dataFilter, bt)
	p.Start()
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < count; i++ {
			if i%2 == 0 {
				a := &tA{
					s: strconv.Itoa(i),
				}
				sinks[0].Push(a)
			} else {
				b := &tB{
					s: strconv.Itoa(i),
				}
				sinks[0].Push(b)
			}
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		for j := 0; j < count/2-1; j++ {
			data := sources[1].Pull()
			tag := data.Tag()
			if tag != at {
				t.Fatal(tag)
			}
		}
		wg.Done()
	}()
	for j := 0; j < count/2-1; j++ {
		data := sources[2].Pull()
		tag := data.Tag()
		if tag != bt {
			t.Fatal(tag)
		}
	}
	wg.Wait()
	t.Log("ok")
}
