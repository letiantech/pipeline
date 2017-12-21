package gpool

type Result struct {
	Res interface{}
	Err error
}

type Job struct {
	p   *Pool
	ch  chan *Result
	req interface{}
}

type Pool struct {
	ch chan *Job
	f  func(req interface{}) (res interface{}, err error)
}

func NewPool(size int, f func(req interface{}) (res interface{}, err error)) *Pool {
	p := &Pool{
		ch: make(chan *Job, size),
		f:  f,
	}
	for i := 0; i < size; i++ {
		go func(p1 *Pool) {

		}(p)
	}
	return p
}

func (p *Pool) NewJob(req interface{}) *Job {
	j := &Job{
		p:   p,
		ch:  make(chan *Result),
		req: req,
	}
	return j
}

func (p *Pool) doJob(j *Job) {
	res := &Result{}
	res.Res, res.Err = p.f(j.req)
	j.ch <- res
}

func (p *Pool) Destroy() {

}
