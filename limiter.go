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

package pipeline

import (
	"sync/atomic"
	"time"
)

const (
	//the slowest limited speed
	MinLimitedSpeed = 0.001
	//the fastest limited speed
	MaxLimitedSpeed = 1000.0
)

//Limiter is used to limit speed of pipeline
//by providing a time duration
type Limiter struct {
	duration int64
	nextTime int64
	speed    float32
}

//create a new limiter and set its speed
func NewLimiter(speed float32) *Limiter {
	l := &Limiter{
		nextTime: time.Now().UnixNano() / int64(time.Millisecond),
	}
	l.SetSpeed(speed)
	return l
}

//update the time duration of limiter
func (l *Limiter) Update() time.Duration {
	duration := atomic.LoadInt64(&l.duration)
	if duration <= 0 {
		return 0
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	nextTime := atomic.LoadInt64(&l.nextTime)
	delta := now - nextTime
	if delta > 0 {
		if delta >= duration {
			nextTime = now + duration/4
		} else {
			nextTime = now + duration - delta
		}
		delta = 0
	} else {
		nextTime = now + duration
		delta = -delta
	}
	atomic.StoreInt64(&l.nextTime, nextTime)
	return time.Duration(delta)
}

//set speed of limiter
func (l *Limiter) SetSpeed(speed float32) {
	if speed < 0 || speed > MaxLimitedSpeed {
		atomic.StoreInt64(&l.duration, 0)
		return
	}
	if speed < MinLimitedSpeed {
		speed = MinLimitedSpeed
	}

	atomic.StoreInt64(&l.duration, int64(time.Second*1000)/int64(1000.0*speed))
}
