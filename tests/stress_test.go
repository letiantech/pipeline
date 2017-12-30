package test

import (
	"fmt"
	"pipeline"
	"testing"

	"reflect"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPipeline is used to test pipeline.PipeLine
func TestPipelineStress(t *testing.T) {
	type TestType2 struct {
		A string
	}
	type TestType1 struct {
		A int
	}
	type TestType3 struct {
		A int
	}
	cfg := &pipeline.Config{
		Name:       "base",
		BufferSize: 100,
		PoolSize:   100,
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
	var task *pipeline.Task
	for i := 0; i < 1000000; i++ {
		data0 := &TestType1{A: i}
		task = p.PushTask(data0)
	}
	end := time.Now().UnixNano()
	var v1 interface{}
	select {
	case v1 = <-task.Chan():
		end = time.Now().UnixNano()
	}
	dt := (end - start) / 1000000
	fmt.Println(dt)
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("v1 Should Be *TestType1", func() {
			So(reflect.TypeOf(v1), ShouldEqual, reflect.TypeOf(&TestType1{}))
		})
	})
	p.DestroyAll()
}
