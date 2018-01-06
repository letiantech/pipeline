package main

import (
	"fmt"
	"pipeline"
	"reflect"
)

func main() {
	type TestType2 struct {
		A string
	}
	type TestType1 struct {
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
			fmt.Println(data)
			d := data.(*TestType1)
			d.A++
			if d.A%2 == 0 {
				return d
			} else {
				return &TestType2{
					A: "abc",
				}
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
					fmt.Println(data)
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
					fmt.Println(data)
					ch2 <- data
					return data
				},
			},
		},
	}
	p, err := pipeline.NewPipeline(cfg)
	if err != nil {
		fmt.Println(err)
	}
	data := &TestType1{A: 1}
	p.Push(data)
	p.Push(data)
	v1 := <-ch1
	v2 := <-ch2
	fmt.Println(v1)
	fmt.Println(v2)
	p.CloseAll()
}
