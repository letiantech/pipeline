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

package main

import (
	"fmt"
	"reflect"

	"github.com/letiantech/pipeline"
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
