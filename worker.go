package gopool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var workerPool sync.Pool

func init() {
	workerPool.New = newWorker
}

var workerid uint32 = 0

type worker struct {
	id           uint32
	q            *taskQueue
	deferHandler func()
}

func newWorker() interface{} {
	id := atomic.AddUint32(&workerid, 1)
	fmt.Printf("new worker:%d\n", id)
	return &worker{id: id}
}

func (w *worker) run() {
	go func() {
		defer func() {
			if w.deferHandler != nil {
				w.deferHandler()
			}
			w.recycle()
		}()
		// run tasks one by one
		for {
			t := w.q.dequeue()
			if t == nil {
				return
			}
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("panic:", e)
				}
			}()
			fmt.Printf("running[wid:%d][tid:%d]\n", w.id, t.id)
			t.f()
			t.recycle()
		}

	}()
}

func (w *worker) reset() {
	w.q = nil
}

func (w *worker) recycle() {
	w.reset()
	workerPool.Put(w)
}
