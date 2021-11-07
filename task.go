package gopool

import (
	"sync"
	"sync/atomic"
)

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

type task struct {
	id   uint32
	f    func()
	next *task
}

func (t *task) reset() {
	t.f = nil
	t.next = nil
}

func (t *task) recycle() {
	t.reset()
	taskPool.Put(t)
}

var taskId uint32 = 0

func newTask() interface{} {
	id := atomic.AddUint32(&taskId, 1)
	return &task{id: id}
}

type taskQueue struct {
	head, tail *task
	m          sync.Mutex
}

func (q *taskQueue) enqueue(t *task) {
	q.m.Lock()
	if q.head == nil {
		q.head = t
		q.tail = t
	} else {
		q.tail.next = t
		q.tail = t
	}
	q.m.Unlock()
}

func (q *taskQueue) dequeue() *task {
	q.m.Lock()
	if q.head == nil {
		return nil
	}
	t := q.head
	q.head = q.head.next
	q.m.Unlock()
	return t
}
