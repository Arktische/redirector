package redirector

import (
	"math/rand"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	demofunc := func(args interface{}) {
		sleeptime := rand.Int63n(500)
		t.Log("task consumes", sleeptime, "ms", ", task id is", args.(int))
		time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	}
	pool, _ := NewPool(20)
	for i := 0; i < 2000; i++ {
		err := pool.Submit(demofunc, i)
		if err != nil {
			t.Log(err)
		}
	}
}

func benchfunc(args interface{}) {
	time.Sleep(time.Duration(rand.Int63n(200)) * time.Millisecond)
}

func BenchmarkPool(b *testing.B) {
}
