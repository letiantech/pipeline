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
	"time"

	"github.com/letiantech/pipeline"
)

type testData int

func (td testData) Tag() string {
	return "int"
}

func TestPipe(t *testing.T) {
	const speedIn = 100
	const speedOut = 50
	const count = 300
	const size = 2
	p := new(pipeline.BasePipe).Init(size)
	minSpeed := speedIn
	if minSpeed > speedOut {
		minSpeed = speedOut
	}
	source := p.GetSource()
	sink := p.GetSink()
	source.SetSpeed(speedOut)
	sink.SetSpeed(speedIn)
	go func() {
		for i := 0; i < count; i++ {
			sink.Push(testData(i))
		}
	}()
	start := time.Now().Unix()
	for j := 0; j < count; j++ {
		source.Pull()
	}
	end := time.Now().Unix()
	if int(end-start+1) < count/minSpeed {
		t.Fatal(end-start, count/minSpeed)
	}
	t.Log(start, end)
}
