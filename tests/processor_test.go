package test

import (
	"testing"

	"gpool"

	. "github.com/smartystreets/goconvey/convey"
)

// TestProcessor is used to test gpool.Processor
func TestProcessor(t *testing.T) {
	out, _ := gpool.NewPipeLine(1)
	in, _ := gpool.NewPipeLine(1)
	cfg := &gpool.ProcessorConfig{
		Speed: 2,
		In:    in,
		Out:   out,
		Func: func(data interface{}) interface{} {
			ret := data.(int32) + 1
			return int32(ret)
		},
	}
	p := gpool.NewProcessor(cfg)
	p.In(int32(1))
	d := p.Out()
	v := int32(0)
	if d != nil {
		v = d.(int32)
	}
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Value Should Be 2", func() {
			So(v, ShouldEqual, 2)
		})
	})
}
