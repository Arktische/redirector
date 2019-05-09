package redirector

import (
	"errors"
)

// LoopQueue is looped-queue
type LoopQueue struct {
	data []interface{}
	head int
	tail int
	size int
}

// NewLoopQueue returns a new looped queue
func NewLoopQueue(size int) *LoopQueue {
	if size <= 0 {
		return nil
	}
	queue := &LoopQueue{
		data: make([]interface{}, size+1),
		size: size,
	}
	return queue
}

// Empty returns whether the queue is empty
func (q *LoopQueue) Empty() bool {
	return q.head == q.tail
}

// Full returns whether the queue is full
func (q *LoopQueue) Full() bool {
	return (q.tail+1)%q.size == q.head
}

// Enqueue insert element into the queue
func (q *LoopQueue) Enqueue(elem interface{}) error {
	if q.size == 0 {
		return errors.New("zero size queue")
	}
	if q.Full() {
		return errors.New("queue is full")
	}
	q.data[q.tail] = elem
	q.tail = (q.tail + 1) % q.size
	return nil
}

// Dequeue pop out element from the queue
func (q *LoopQueue) Dequeue() interface{} {
	if q.Empty() {
		return nil
	}
	res := q.data[q.head]
	q.head = (q.head + 1) % q.size
	return res
}
