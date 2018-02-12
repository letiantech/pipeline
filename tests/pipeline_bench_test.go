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
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/letiantech/pipeline"
)

// BenchmarkPipeline is used to run benchmark test for pipeline.PipeLine
func BenchmarkPipeline(b *testing.B) {
	type TestType struct {
		A int
	}
	cfg := &pipeline.Config{
		Name:       "base",
		BufferSize: 100,
		PoolSize:   runtime.GOMAXPROCS(0),
		Type:       reflect.TypeOf(&TestType{}),
		Func: func(data interface{}) interface{} {
			return data
		},
	}
	p, err := pipeline.NewPipeline(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	b.ResetTimer()
	var task pipeline.Task
	for i := 0; i < b.N; i++ {
		data0 := &TestType{A: i}
		task = pipeline.NewTask(data0, true)
		p.Push(task)
	}
	<-task.Chan()
	p.CloseAll()
}
