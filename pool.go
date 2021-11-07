package gopool

import (
	"sync/atomic"
)

type Pool interface {
	Go(func())
	Close()
}

func NewPool(cap int32) Pool {
	return &pool{
		cap:       cap,
		workerNum: 0,
		q:         &taskQueue{},
	}
}

type pool struct {
	cap, workerNum int32
	q              *taskQueue
	closed         int32
}

func (p *pool) Go(f func()) {
	if atomic.LoadInt32(&p.closed) == 1 {
		panic("pool closed")
	}
	t := taskPool.Get().(*task)
	t.f = f
	p.q.enqueue(t) // enqueue, wait for available worker to run
	if atomic.AddInt32(&p.workerNum, 1) <= p.cap {
		// get woker, and run
		w := workerPool.Get().(*worker)
		w.q = p.q
		w.deferHandler = func() {
			atomic.AddInt32(&p.workerNum, -1) //dec worknum
		}
		w.run()
	}
}

func (p *pool) Close() {
	atomic.StoreInt32(&p.closed, 1)
}
