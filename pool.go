package gpool

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Pool struct {
	capacity int32
	running  int32

	lock    *sync.Mutex
	workers workerQueue
	cond    *sync.Cond

	cache   *sync.Pool
	options Options
}

var (
	ErrPoolOverload = errors.New("gPool overload")
)

func New(options Options) *Pool {
	pool := &Pool{
		capacity: max(-1, options.poolCapacity),
		lock:     &sync.Mutex{},
		workers:  newWorkerQueue(options.queueType, options.poolCapacity),
		options:  options,
		cache:    &sync.Pool{},
	}
	pool.cache.New = func() any {
		return &worker{
			pool:  pool,
			tasks: make(chan func(), min(minTaskCapacity, options.taskCapacity)),
		}
	}
	pool.cond = sync.NewCond(pool.lock)
	return pool
}

func (p *Pool) Submit(task func()) error {
	// no available worker and pool is full
	if w := p.getWorker(true); w == nil {
		return ErrPoolOverload
	} else {
		w.tasks <- task
	}
	return nil
}

func (p *Pool) getWorker(blocking bool) *worker {
	p.lock.Lock()
	// poll worker from queue
	if w := p.workers.poll(); w != nil {
		p.lock.Unlock()
		return w
	} else if c := p.Capacity(); c != -1 && c > p.Running() {
		// create a new worker if pool is not full
		w := p.cache.Get().(*worker)
		w.run()
		p.lock.Unlock()
		return w
	} else {
		// no available worker and pool is full, return nil in non-blocking mode
		if !blocking {
			p.lock.Unlock()
			return nil
		}
	}
RETRY:
	// wait for other task to release a worker
	p.cond.Wait()
	if w := p.workers.poll(); w != nil {
		p.lock.Unlock()
		return w
	} else {
		goto RETRY
	}
}

func (p *Pool) returnWorker(w *worker) {
	p.lock.Lock()
	p.workers.put(w)
	p.cond.Signal()
	p.lock.Unlock()
}

func (p *Pool) Capacity() int {
	return int(p.capacity)
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) addRunning(delta int32) {
	atomic.AddInt32(&p.running, delta)
}

func max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
