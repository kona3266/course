package controller

import (
	"sync"
	"sync/atomic"
)

type refCounter struct {
	counter int32
	lock    *sync.Mutex
}

type MultipleLock interface {
	Lock(interface{})
	Unlock(interface{})
}

type lock struct {
	inUse sync.Map
	pool  *sync.Pool
}

func (l *lock) Lock(key interface{}) {
	m := l.getLocker(key)
	atomic.AddInt32(&m.counter, 1)
	m.lock.Lock()
}

func (l *lock) getLocker(key interface{}) *refCounter {
	r, _ := l.inUse.LoadOrStore(key, &refCounter{counter: 0, lock: l.pool.Get().(*sync.Mutex)})
	return r.(*refCounter)
}

func (l *lock) Unlock(key interface{}) {
	m := l.getLocker(key)
	m.lock.Unlock()
	atomic.AddInt32(&m.counter, -1)
	if m.counter <= 0 {
		l.pool.Put(m.lock)
		l.inUse.Delete(key)
	}
}

func NewMultipleLock() MultipleLock {
	return &lock{pool: &sync.Pool{New: func() interface{} {
		return &sync.Mutex{}
	}}}
}
