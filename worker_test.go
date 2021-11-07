package gopool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerBase(t *testing.T) {
	var num int32 = 0
	f := func() {
		atomic.AddInt32(&num, 1)
	}

	q := makeTaskQuene(f, 4)

	w := workerPool.Get().(*worker)
	w.q = q
	w.run()

	<-time.After(time.Millisecond* 200)
}
