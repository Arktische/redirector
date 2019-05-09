package redirector

import (
	"errors"
	"sync"
	"sync/atomic"
)

type f func(interface{})

var (
	// ErrInvalidPoolSize will be returned when setting a negative number as pool capacity, this error will be only used
	// by pool with func because pool without func can be infinite by setting up a negative capacity.
	ErrInvalidPoolSize = errors.New("invalid size for pool")

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")
)

// PoolConfig is config of the goroutine pool
type PoolConfig struct {
	// maxblocked limits maximum blocked goroutine on cond, Once exceed
	// Submit() will fail
	maxblocked int32
	blocking   bool
}

func defaultPoolConfig() (config *PoolConfig) {
	return &PoolConfig{
		maxblocked: 2000,
		blocking:   true,
	}
}

// GoroutinePool goroutine pool struct
type GoroutinePool struct {
	// capacity is goroutine capacity
	capacity int32
	// nrunning counts current goroutines on running state
	nrunning int32
	// nblocked counts goroutines on blocked state
	nblocked int32

	// workers is a looped-queue stores workers
	workers *LoopQueue
	// objCache is a object pool reducing GC pressure
	objCache sync.Pool
	// lock is a spinlock for sync ops
	lock sync.Locker
	// cond is a conditional variable for getting idle worker
	cond *sync.Cond

	config *PoolConfig
}

func (p *GoroutinePool) getWorker() (w *Worker) {
	generateWorker := func() {
		w = p.objCache.Get().(*Worker)
		w.run()
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	var ok bool
	w, ok = p.workers.Dequeue().(*Worker)
	// workers queue have available workers
	if ok {
		return
	}
	if p.nrunning < p.capacity {
		// workers queue is empty but numbers of goroutines
		// have not reached p.capacity yet
		generateWorker()
	} else {
		if !p.config.blocking {
			return
		}
		// needs to wait for workers queue not empty
		for p.workers.Empty() {
			if p.nblocked > p.config.maxblocked {
				return
			}
			p.nblocked++
			p.cond.Wait()
			p.nblocked--
			//
			if p.nrunning == 0 {
				generateWorker()
				return
			}
		}
		w = p.workers.Dequeue().(*Worker)
	}
	return
}

func (p *GoroutinePool) putWorker(worker *Worker) bool {
	p.lock.Lock()
	err := p.workers.Enqueue(worker)
	if err != nil {
		p.lock.Unlock()
		return false
	}
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// Submit submits a task to this pool
func (p *GoroutinePool) Submit(task f, args interface{}) error {
	w := p.getWorker()
	if w == nil {
		return ErrPoolOverload
	}
	w.sendTask(task)
	w.sendArgs(args)
	return nil
}

// AsyncSubmit submits a task asynchronously
func (p *GoroutinePool) AsyncSubmit(task f, args interface{}, err chan error) {
	go func() {
		err <- p.Submit(task, args)
	}()
	return
}

// Resize changes the capacity of pool
func (p *GoroutinePool) Resize(size int32) {
	if size <= 0 || size == p.capacity {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
}

// NewPool returns a new default goroutine pool
func NewPool(size int) (*GoroutinePool, error) {
	p := &GoroutinePool{
		capacity: int32(size),
		lock:     NewSpinLock(),
		config:   defaultPoolConfig(),
		workers:  NewLoopQueue(size),
	}
	p.cond = sync.NewCond(p.lock)
	p.objCache.New = func() interface{} {
		return &Worker{
			pool: p,
			task: make(chan f, 1),
			args: make(chan interface{}, 1),
		}
	}
	return p, nil
}

// NewPoolWithConfig returns a new goroutine pool with customized config
func NewPoolWithConfig(size int, conf *PoolConfig) (*GoroutinePool, error) {
	p := &GoroutinePool{
		capacity: int32(size),
		lock:     NewSpinLock(),
		config:   conf,
		workers:  NewLoopQueue(size),
	}
	p.cond = sync.NewCond(p.lock)
	p.objCache.New = func() interface{} {
		return &Worker{
			pool: p,
			task: make(chan f, 1),
			args: make(chan interface{}, 1),
		}
	}
	return p, nil
}
