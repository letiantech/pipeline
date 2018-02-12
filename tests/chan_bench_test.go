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
	"runtime"
	"sync"
	"testing"

	"github.com/letiantech/pipeline"
)

func benchTestConcurrently(count int64, exec func(j int64), finish func()) {
	wg := sync.WaitGroup{}
	max := (runtime.GOMAXPROCS(0) + 1) / 2
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for k := int64(0); k < count; k += int64(max) {
				exec(k)
			}
		}(i)
	}
	wg.Wait()
	finish()
}

func benchTest(count int64, exec func(j int64), finish func()) {
	for i := int64(0); i < count; i++ {
		exec(i)
	}
	finish()
}

func initChan() chan interface{} {
	ncpus := runtime.GOMAXPROCS(0)
	ch := make(chan interface{}, ncpus)
	max := (ncpus + 2) / 3
	for i := 0; i < max; i++ {
		go func() {
			for {
				d1 := <-ch
				if d1 == nil {
					return
				}
			}
		}()
	}
	return ch
}

func initSyncChan() *pipeline.Chan {
	ncpus := runtime.GOMAXPROCS(0)
	ch := pipeline.MakeChan(ncpus)
	max := (ncpus + 2) / 3
	for i := 0; i < max; i++ {
		go func() {
			for {
				d1 := <-ch.Chan()
				if d1 == nil {
					return
				}
			}
		}()
	}
	return ch
}

func BenchmarkChan(t *testing.B) {
	ch := initChan()
	t.ResetTimer()
	benchTest(int64(t.N), func(j int64) {
		ch <- j
	}, func() {
		close(ch)
	})
}

func BenchmarkSyncChan(t *testing.B) {
	ch := initSyncChan()
	t.ResetTimer()
	benchTest(int64(t.N), func(j int64) {
		ch.Push(j)
	}, func() {
		ch.Close()
	})
}

func BenchmarkChanConcurrently(t *testing.B) {
	ch := initChan()
	t.ResetTimer()
	benchTestConcurrently(int64(t.N), func(j int64) {
		ch <- j
	}, func() {
		close(ch)
	})
}

func BenchmarkSyncChanConcurrently(t *testing.B) {
	ch := initSyncChan()
	t.ResetTimer()
	benchTestConcurrently(int64(t.N), func(j int64) {
		ch.Push(j)
	}, func() {
		ch.Close()
	})
}
