package redirector

import (
	"fmt"
	"testing"
)

type TestType struct {
	a int
	b int
}

func TestLoopQueue(t *testing.T) {
	testStruct := &TestType{
		a: 10,
		b: 11,
	}
	q := NewLoopQueue(10)
	q.Enqueue(testStruct)
	q.Enqueue(testStruct)
	fmt.Println(q.Empty())
	q.Dequeue()
}
