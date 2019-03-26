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

type testTaskData struct {
	A int
}

func (tt *testTaskData) Tag() string {
	return "testTask"
}

func TestTask(t *testing.T) {
	data := &testTaskData{
		A: 1,
	}
	task := pipeline.NewTask(data, 100*time.Millisecond)
	go func() {
		time.Sleep(100 * time.Millisecond)
		task.Finish()
	}()
	select {
	case _ = <-task.Finished():
		t.Log("finished")
	case _ = <-task.Canceled():
		t.Log("canceled")
	}
	data1 := task.GetData()
	if _, ok := data1.(*testTaskData); !ok {
		t.Fatal("failed")
	}
	t.Log("ok")
}
