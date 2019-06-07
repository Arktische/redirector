package redirector

import (
	"math/rand"
	"sort"
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

func TestHashTrieTree(t *testing.T) {
	r := NewHashTrieTree()
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
	testvecsize   = 512
	inputvecsize  = 2000
	teststringlen = 256
)

var testvec [testvecsize]string
var inputvec [inputvecsize]string
var flag bool
var r *RadixNode = NewRadixTree()
var m map[string]interface{} = make(map[string]interface{})
var h *HashTrieNode = NewHashTrieTree()

func readtestvec() {
	for i := 0; i < testvecsize; i++ {
		testvec[i] = getRandomString(teststringlen)
	}
	sort.Strings(testvec[:])
	flag = true
	for i := 0; i < testvecsize; i++ {
		r.AddRoute(testvec[i], nil)
		h.AddRoute(testvec[i], nil)
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
func BenchmarkHashTrie(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	if !flag {
		readinputvec()
		readinputvec()
	}
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		_, _ = h.FindRoute(inputvec[j%inputvecsize])
	}
}
func BenchmarkDoubleArray(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	builder := DoubleArrayTrieBuilder{}
	builder.Build(testvec[:])
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		_, _ = builder.ExactMatchSearch(inputvec[j%inputvecsize])
	}
}
