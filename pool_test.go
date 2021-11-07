package gopool

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoolBase(t *testing.T) {
	p := NewPool(4)

	var num int32 = 0
	f := func() {
		time.Sleep(time.Millisecond * 10)
		atomic.AddInt32(&num, 1)
	}

	const c = 10
	for i := 0; i < c; i++ {
		p.Go(f)
	}

	<-time.After(time.Millisecond * 200)

	assert.Equal(t, int32(c), num)

}

func TestPoolPanic(t *testing.T) {
	p := NewPool(2)

	p.Go(func() {
		panic("as expected")
	})

	<-time.After(time.Millisecond * 50)
}
