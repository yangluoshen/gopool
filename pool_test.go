package gopool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoolBase(t *testing.T) {

	var num int32 = 0
	var wg sync.WaitGroup
	const c = 100
	wg.Add(c)
	f := func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 10)
		atomic.AddInt32(&num, 1)
	}

	p := NewPool(4)
	for i := 0; i < c; i++ {
		p.Go(f)
	}

	wg.Wait()

	assert.Equal(t, int32(c), num)
}

func TestPoolPanic(t *testing.T) {
	p := NewPool(2)

	p.Go(func() {
		panic("as expected")
	})

	<-time.After(time.Millisecond * 50)
}
