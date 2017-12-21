package test

import (
	"gpool"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPipeline is used to test gpool.PipeLine
func TestPipeline(t *testing.T) {
	p, _ := gpool.NewPipeLine(1)
	p.Start()
	p.In(int32(1))
	d := p.Out()
	v := int32(0)
	if d != nil {
		v = d.(int32)
	}
	p.Close()
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Value Should Be 1", func() {
			So(v, ShouldEqual, 1)
		})
	})
}
