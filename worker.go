package redirector

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// Worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type Worker struct {
	pool *GoroutinePool
	// task is a job should be done.
	task chan f
	args chan interface{}
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *Worker) run() {
	atomic.AddInt32(&w.pool.nrunning, 1)
	go func() {
		defer func() {
			atomic.AddInt32(&w.pool.nrunning, -1)
			w.pool.objCache.Put(w)
			p := recover()
			if p != nil {
				fmt.Println("worker exits from panic:", p)
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				fmt.Printf("stack frame is\n%s\n", string(buf[:n]))
			}
		}()
		//
		for f := range w.task {
			if f == nil {
				return
			}
			f(<-w.args)

			// collect and reuse
			if ok := w.pool.putWorker(w); !ok {
				return
			}
		}
	}()
}

// stop this worker.
func (w *Worker) stop() {
	w.sendTask(nil)
}

// sendTask sends a task to this worker.
func (w *Worker) sendTask(task f) {
	w.task <- task
}

func (w *Worker) sendArgs(args interface{}) {
	w.args <- args
}
