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
	"math/rand"
	"testing"

	"github.com/letiantech/pipeline"
)

func TestMakeChan(t *testing.T) {
	testChan := pipeline.MakeChan(2)
	if testChan == nil {
		t.Fatal("testChan is nil")
	}
	testChan.Close()
}

func TestChanPushGet(t *testing.T) {
	testChan := pipeline.MakeChan(2)
	var r = rand.New(rand.NewSource(1))
	data := r.Int63()
	err := testChan.Push(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	d := testChan.Pull()
	if d == nil {
		t.Fatal("data is nil")
	}
	switch v := d.(type) {
	case int64:
		if v != data {
			t.Fatal("data value not match")
		}
		break
	default:
		t.Fatal("data type not match")
	}
	testChan.Close()
}

// test push after channel is closed
func TestChanClosePush(t *testing.T) {
	var r = rand.New(rand.NewSource(1))
	testChan := pipeline.MakeChan(2)
	data1 := r.Int63()
	ch := testChan.Chan()
	err := testChan.Push(data1)
	if err != nil {
		t.Fatal(err.Error())
	}
	testChan.Close()
	if !testChan.IsClosed() {
		t.Fatal("sync channel is not closed")
	}
	data2 := r.Int63()
	err = testChan.Push(data2)
	if err == nil {
		t.Fatal("sync channel is not closed")
	}
	switch v := testChan.Pull().(type) {
	case int64:
		if v != data1 {
			t.Fatal("data value not match")
		}
		break
	default:
		t.Fatal("data type not match")
	}
	v := <-ch
	if v != nil {
		t.Fatal("channel is not closed")
	}
	testChan = nil
}
