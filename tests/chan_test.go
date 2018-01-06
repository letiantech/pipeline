package test

import (
	"math/rand"
	"pipeline"
	"testing"
)

func TestMakeChan(t *testing.T) {
	testChan := pipeline.MakeChan(2)
	if testChan == nil {
		t.Fatal("testChan is nil")
	}
	testChan.Close()
}

func TestChanPushGet(t *testing.T) {
	testChan := pipeline.MakeChan(2)
	var r = rand.New(rand.NewSource(1))
	data := r.Int63()
	err := testChan.Push(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	d := testChan.Get()
	if d == nil {
		t.Fatal("data is nil")
	}
	switch v := d.(type) {
	case int64:
		if v != data {
			t.Fatal("data value not match")
		}
		break
	default:
		t.Fatal("data type not match")
	}
	testChan.Close()
}

// test push after channel is closed
func TestChanClosePush(t *testing.T) {
	var r = rand.New(rand.NewSource(1))
	testChan := pipeline.MakeChan(2)
	data1 := r.Int63()
	ch := testChan.Chan()
	err := testChan.Push(data1)
	if err != nil {
		t.Fatal(err.Error())
	}
	testChan.Close()
	if !testChan.IsClosed() {
		t.Fatal("sync channel is not closed")
	}
	data2 := r.Int63()
	err = testChan.Push(data2)
	if err == nil {
		t.Fatal("sync channel is not closed")
	}
	switch v := testChan.Get().(type) {
	case int64:
		if v != data1 {
			t.Fatal("data value not match")
		}
		break
	default:
		t.Fatal("data type not match")
	}
	v := <-ch
	if v != nil {
		t.Fatal("channel is not closed")
	}
	testChan = nil
}
