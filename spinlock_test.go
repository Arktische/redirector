package redirector

import (
	"sync"
	"testing"
)

func TestSpinlock(t *testing.T) {
	lock := NewSpinLock()
	var i int = 0
	var sum int = 0
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for true {
			lock.Lock()
			if i > 100 {
				lock.Unlock()
				break
			}
			sum += i
			i++
			t.Log(i)
			lock.Unlock()
		}
	}()
	for true {
		lock.Lock()
		if i > 100 {
			lock.Unlock()
			break
		}
		sum += i
		i++
		t.Log(i)
		lock.Unlock()
	}

	wg.Wait()
	t.Log("result is", sum)
	if sum != 5050 {
		t.Failed()
	}
}
