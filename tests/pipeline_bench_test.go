package test

import (
	"fmt"
	"pipeline"
	"reflect"
	"runtime"
	"testing"
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
