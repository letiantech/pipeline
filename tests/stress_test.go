package test

import (
	"pipeline"
	"testing"

	"reflect"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPipeline is used to test pipeline.PipeLine
func TestPipelineStress(t *testing.T) {
	type TestType1 struct {
		A int
	}
	cfg := &pipeline.Config{
		Name:       "base",
		BufferSize: 1000,
		PoolSize:   1000,
		Type:       reflect.TypeOf(&TestType1{}),
		Func: func(data interface{}) interface{} {
			time.Sleep(10 * time.Millisecond)
			return data
		},
	}
	p, err := pipeline.NewPipeline(cfg)
	if err != nil {
		Convey("Subject: Test Station Endpoint\n", t, func() {
			Convey("Value Should Be nil", func() {
				So(err, ShouldEqual, nil)
			})
		})
	}
	start := time.Now().UnixNano()
	var task pipeline.Task
	for i := 0; i < 1000000; i++ {
		data0 := &TestType1{A: i}
		task = pipeline.NewTask(data0, true)
		p.Push(task)
	}
	end := time.Now().UnixNano()
	select {
	case <-task.Chan():
		end = time.Now().UnixNano()
	}
	dt := (end - start) / 1000000000
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("v1 Should Be *TestType1", func() {
			So(dt, ShouldEqual, 10)
		})
	})
	p.CloseAll()
}
