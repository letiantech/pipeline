package test

import (
	"pipeline"
	"runtime"
	"sync"
	"testing"
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
