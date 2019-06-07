package redirector

import (
	"math/rand"
	"testing"
	"time"
)

const (
	testSize = 256
	distBase = 4
	distNum  = 1
)

func TestHashMap(t *testing.T) {
	var kvec [testSize]int
	var vvec [testSize]int
	m := NewHashMap()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < testSize; i++ {
		kvec[i] = rand.Intn(testSize)
		vvec[i] = rand.Intn(testSize)
		m.Set(kvec[i], vvec[i])
	}
	var inpvec [testSize]int
	var outvec [testSize]bool

	for i := 0; i < testSize; i++ {
		p := rand.Intn(distBase)
		if p <= distNum {
			inpvec[i] = rand.Intn(testSize) + testSize
			outvec[i] = false
		} else {
			inpvec[i] = kvec[i]
			outvec[i] = true
		}
	}

	for i := 0; i < testSize; i++ {
		val, ok := m.Get(inpvec[i])
		t.Log("key:", inpvec[i], "get val:", val, ok)
		if ok != outvec[i] {
			t.Fail()
		}
	}
}

func BenchmarkLinkedAddressHashMap(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	var kvec [testSize]int
	var vvec [testSize]int
	m := NewHashMap()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < testSize; i++ {
		kvec[i] = rand.Intn(testSize)
		vvec[i] = rand.Intn(testSize)
		m.Set(kvec[i], vvec[i])
	}
	var inpvec [testSize]int
	var outvec [testSize]bool

	for i := 0; i < testSize; i++ {
		p := rand.Intn(distBase)
		if p <= distNum {
			inpvec[i] = rand.Intn(testSize) + testSize
			outvec[i] = false
		} else {
			inpvec[i] = kvec[i]
			outvec[i] = true
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < testSize; i++ {
			m.Get(inpvec[i])
		}
	}
}

func BenchmarkOriginalMap(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	var kvec [testSize]int
	var vvec [testSize]int
	m := make(map[int]int, 8)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < testSize; i++ {
		kvec[i] = rand.Intn(testSize)
		vvec[i] = rand.Intn(testSize)
		m[kvec[i]] = vvec[i]
	}
	var inpvec [testSize]int
	var outvec [testSize]bool

	for i := 0; i < testSize; i++ {
		p := rand.Intn(distBase)
		if p <= distNum {
			inpvec[i] = rand.Intn(testSize) + testSize
			outvec[i] = false
		} else {
			inpvec[i] = kvec[i]
			outvec[i] = true
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < testSize; i++ {
			_, _ = m[inpvec[i]]
		}
	}
}
