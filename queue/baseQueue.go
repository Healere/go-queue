package queue

import (
	"errors"
	"sync"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueFull  = errors.New("queue is full")
)

// Queue 基础队列接口
type Queue interface {
	Enqueue(interface{}) error
	Dequeue() (interface{}, error)
	Peek() (interface{}, error)
	Size() int
	IsEmpty() bool
	IsFull() bool
	Clear()
}

// SliceQueue 基于切片的队列实现
type SliceQueue struct {
	items    []interface{}
	capacity int
	mu       sync.RWMutex
}

// NewQueue 创建新的切片队列
func NewQueue(capacity int) *SliceQueue {
	return &SliceQueue{
		items:    make([]interface{}, 0, capacity),
		capacity: capacity,
	}
}

func (q *SliceQueue) Enqueue(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return ErrQueueFull
	}
	q.items = append(q.items, item)
	return nil
}

func (q *SliceQueue) Dequeue() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *SliceQueue) Peek() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return q.items[0], nil
}

func (q *SliceQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

func (q *SliceQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items) == 0
}

func (q *SliceQueue) IsFull() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.capacity > 0 && len(q.items) == q.capacity
}

func (q *SliceQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = make([]interface{}, 0, q.capacity)
}
