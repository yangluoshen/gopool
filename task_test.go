package gopool

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTaskQuene(f func(), n int) *taskQueue {
	q := &taskQueue{}
	for i:=0; i<n; i++ {
		t := taskPool.Get().(*task)
		t.f = f
		q.enqueue(t)
	}
	return q
}

func TestTaskBase(t *testing.T) {
	var num int32 = 0
	f := func() {
		atomic.AddInt32(&num, 1)
	}

	q := makeTaskQuene(f, 2) 

	var i uint32 = 1
	for tp := q.head; tp != nil; tp = tp.next {
		tp.f()
		assert.Equal(t, uint32(i), tp.id)
		i ++
		t.Logf("tid:%d, num:%d", tp.id, num)
	}

	tx := q.dequeue()
	assert.Equal(t, uint32(1), tx.id)

	for tp := q.head; tp != nil; tp = tp.next {
		t.Logf("tid:%d, num:%d", tp.id, num)
	}
}
