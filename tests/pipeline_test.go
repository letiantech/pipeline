package test

import (
	"pipeline"
	"testing"

	"reflect"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPipeline is used to test pipeline.PipeLine
func TestPipeline(t *testing.T) {
	type TestType2 struct {
		A string
	}
	type TestType1 struct {
		A int
	}
	type TestType3 struct {
		A int
	}
	ch1 := make(chan interface{}, 1)
	ch2 := make(chan interface{}, 1)
	cfg := &pipeline.Config{
		Name:       "base",
		Speed:      2,
		BufferSize: 2,
		PoolSize:   2,
		Type:       reflect.TypeOf(&TestType1{}),
		Func: func(data interface{}) interface{} {
			d := data.(*TestType1)
			switch d.A % 3 {
			case 0:
				return d
			case 1:
				return &TestType2{
					A: "abc",
				}
			}
			return &TestType3{
				A: d.A + 1,
			}
		},
		NextConfigs: []*pipeline.Config{
			{
				Name:       "base",
				Speed:      2,
				BufferSize: 2,
				PoolSize:   1,
				Type:       reflect.TypeOf(&TestType1{}),
				Func: func(data interface{}) interface{} {
					ch1 <- data
					return data
				},
			},
			{
				Name:       "base",
				Speed:      2,
				BufferSize: 2,
				PoolSize:   1,
				Type:       reflect.TypeOf(&TestType2{}),
				Func: func(data interface{}) interface{} {
					ch2 <- data
					return data
				},
			},
			{
				Name:       "base",
				Speed:      2,
				BufferSize: 2,
				PoolSize:   1,
				Type:       reflect.TypeOf(&TestType3{}),
				Func: func(data interface{}) interface{} {
					t := data.(*TestType3)
					t.A = 3
					return data
				},
			},
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
	data0 := &TestType1{A: 0}
	p.Push(data0)
	data1 := &TestType1{A: 1}
	p.Push(data1)
	data2 := &TestType1{A: 2}
	dt := pipeline.NewTask(data2, true)
	p.Push(dt)
	v1 := <-ch1
	v2 := <-ch2
	v3 := <-dt.Chan()
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("v1 Should Be *TestType1", func() {
			So(reflect.TypeOf(v1), ShouldEqual, reflect.TypeOf(&TestType1{}))
		})
		Convey("v2 Should Be *TestType2", func() {
			So(reflect.TypeOf(v2), ShouldEqual, reflect.TypeOf(&TestType2{}))
		})
		Convey("v3 Should Be *TestType3", func() {
			So(reflect.TypeOf(v3), ShouldEqual, reflect.TypeOf(&TestType3{}))
		})
	})
	p.CloseAll()
}
