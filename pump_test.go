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
	"testing"

	"reflect"

	"strconv"

	"sync"

	"github.com/letiantech/pipeline"
)

func TestPump(t *testing.T) {
	pipe0 := new(pipeline.BasePipe).Init(10)
	pipe1 := new(pipeline.BasePipe).Init(10)
	pipe2 := new(pipeline.BasePipe).Init(10)
	source0 := pipe0.GetSource()
	source1 := pipe1.GetSource()
	source2 := pipe2.GetSource()
	sink0 := pipe0.GetSink()
	sink1 := pipe1.GetSink()
	sink2 := pipe2.GetSink()
	p := new(pipeline.BasePump).Init(source0, func(data interface{}) (interface{}, string) {
		return data, reflect.TypeOf(data).String()
	})
	const speedIn = 100
	const speedOut = 50
	const count = 300
	var typeOfInt = reflect.TypeOf(int(1)).String()
	var typeOfStr = reflect.TypeOf("").String()
	minSpeed := speedIn
	if minSpeed > speedOut {
		minSpeed = speedOut
	}
	source0.SetSpeed(speedOut)
	source1.SetSpeed(speedOut)
	source2.SetSpeed(speedOut)
	sink0.SetSpeed(speedIn)
	sink1.SetSpeed(speedIn)
	sink2.SetSpeed(speedIn)
	p.BindSink(sink1, typeOfInt)
	p.BindSink(sink2, typeOfStr)
	p.Start()
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < count/2; i++ {
			sink0.Push(int(i))
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		for i := 0; i < count/2; i++ {
			sink0.Push(strconv.Itoa(i))
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		for j := 0; j < count/2; j++ {
			tp := reflect.TypeOf(source1.Pull()).String()
			if tp != typeOfInt {
				t.Fatal(tp)
			}
		}
		wg.Done()
	}()
	for j := 0; j < count/2; j++ {
		tp := reflect.TypeOf(source2.Pull()).String()
		if tp != typeOfStr {
			t.Fatal(tp)
		}
	}
	wg.Wait()
	t.Log("ok")
}
