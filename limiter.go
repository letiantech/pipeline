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
	duration  int64
	nextTime  int64
	speed     float32
	unlimited int32
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
