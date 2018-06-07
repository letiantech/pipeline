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

package test

import (
	"testing"

	"time"

	"github.com/letiantech/pipeline"
)

type TestType struct {
	A int
}

func (tt *TestType) Tag() string {
	return "TestType"
}

var ch = make(chan pipeline.Data, 1000)

func filter(data pipeline.Data) pipeline.Data {
	if task, ok := data.(pipeline.Task); ok {
		ch <- task.GetData()
		task.Finish()
	}
	return nil
}

// BenchmarkPump is used to run benchmark test for pipeline.Pump
func BenchmarkPump(b *testing.B) {
	go func() {
		tag := (&TestType{}).Tag()
		pipe := new(pipeline.BasePipe).Init(100)
		p := new(pipeline.BasePump).Init(pipe.GetSource())
		p.AddFilter(filter, tag)
		p.Start()
		b.ResetTimer()
		var task pipeline.Task
		for i := 0; i < b.N; i++ {
			data0 := &TestType{A: i}
			task = pipeline.NewTask(data0, time.Second)
			pipe.Push(task)
		}
	}()
	for i := 0; i < b.N; i++ {
		_ = <-ch
	}
}
