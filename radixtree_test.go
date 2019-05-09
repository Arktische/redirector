package redirector

import (
	"math/rand"
	"testing"
	"time"
)

func TestRadixTree(t *testing.T) {
	r := NewRadixTree()
	r.AddRoute("192.168.9.1", nil)
	r.AddRoute("192.168.3.1", nil)
	r.AddRoute("34.42.56.7", nil)
	r.AddRoute("CDCD:910A:2222:5498:8475:1111:3900:2020", nil)
	r.AddRoute("3f4c:8a70:5ef2:3425:8392:2912:a34f:8d32", nil)

	t.Log(r.FindRoute("192.168.9.1"))
	t.Log(r.FindRoute("192.145.3.5"))
	t.Log(r.FindRoute("34.42.56.7"))
	t.Log(r.FindRoute("3f4c:8a70:5ef2:3425:8392:2912:a34f:8d32"))

	r.AddRoute("baidu.com", nil)
	r.AddRoute("bilibili.com", nil)
	r.AddRoute("bilibili.tv", nil)
	r.AddRoute("bilibiligame.com", nil)

	t.Log(r.FindRoute("www.bilibili.com"))
	t.Log(r.FindRoute("bilibili.tv"))
}

func TestImap(t *testing.T) {
	m := newimap()
	rand.Seed(2583953)
	var testvec [32]int
	for i := 0; i < 32; i++ {
		testvec[i] = rand.Intn(32) + 1
	}
	for i := 0; i < 32; i++ {
		m.set(testvec[i], i)
	}

	for i := 0; i < 32; i++ {
		t.Log(m.get(testvec[i]))
	}
	m.set(testvec[7], 7)
	t.Log(m.get(testvec[7]))
}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz./?&*"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(38197953)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := rand.Intn(l)
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const (
	testvecsize   = 40960
	inputvecsize  = 20
	teststringlen = 256
)

var testvec [testvecsize]string
var inputvec [inputvecsize]string
var flag bool
var r *RadixNode = NewRadixTree()
var m map[string]interface{} = make(map[string]interface{})

func readtestvec() {
	for i := 0; i < testvecsize; i++ {
		testvec[i] = getRandomString(teststringlen)
	}
	flag = true
	for i := 0; i < testvecsize; i++ {
		r.AddRoute(testvec[i], nil)
		m[testvec[i]] = i
	}
}

func readinputvec() {
	for i := 0; i < inputvecsize; i++ {
		if rand.Int31n(40000000) <= 10000000 {
			inputvec[i] = testvec[rand.Intn(testvecsize)]
		} else {
			inputvec[i] = getRandomString(teststringlen)
		}
	}
	flag = true
}

func BenchmarkHashMap(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	if !flag {
		readtestvec()
		readinputvec()
	}
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		_, _ = m[inputvec[j%inputvecsize]]
	}
}
func BenchmarkRadixTree(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	if !flag {
		readtestvec()
		readinputvec()
	}
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		_, _ = r.FindRoute(inputvec[j%inputvecsize])
	}
}

func BenchmarkDoubleArray(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	builder := DoubleArrayBuilder{}
	builder.Build(testvec[:])
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		_, _ = builder.ExactMatchSearch(inputvec[j%inputvecsize])
	}
}

const (
	threshold = 16
)

func BenchmarkMapVsTraverse_Map(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	m := make(map[int]int)
	for i := 0; i < threshold; i++ {
		m[i] = rand.Int()
	}
	b.StartTimer()
	target := 135
	for j := 0; j < b.N; j++ {
		_, ok := m[target]
		if ok {
			target++
		}
	}
}

func BenchmarkHashMapVsTraverse_Traverse(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	var m [threshold]int
	for i := 0; i < threshold; i++ {
		m[i] = rand.Int()
	}
	b.StartTimer()
	target := 129
	for j := 0; j < b.N; j++ {
		for k := 0; k < len(m); k++ {
			if target == m[k] {
				target++
				break
			}
		}
	}
}
