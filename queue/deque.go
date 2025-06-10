package queue
import (
	"sync"
)

// Deque 双端队列接口
type Deque interface {
	Queue
	EnqueueFront(interface{}) error
	DequeueBack() (interface{}, error)
	PeekFront() (interface{}, error)
	PeekBack() (interface{}, error)
}

// SliceDeque 基于切片的双端队列实现
type SliceDeque struct {
	items    []interface{}
	capacity int
	mu       sync.RWMutex
}

// NewSliceDeque 创建新的双端队列
func NewSliceDeque(capacity int) *SliceDeque {
	return &SliceDeque{
		items:    make([]interface{}, 0, capacity),
		capacity: capacity,
	}
}

func (q *SliceDeque) Enqueue(item interface{}) error {
	return q.EnqueueBack(item)
}

func (q *SliceDeque) Dequeue() (interface{}, error) {
	return q.DequeueFront()
}

func (q *SliceDeque) Peek() (interface{}, error) {
	return q.PeekFront()
}

func (q *SliceDeque) EnqueueFront(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return ErrQueueFull
	}
	q.items = append([]interface{}{item}, q.items...)
	return nil
}

func (q *SliceDeque) EnqueueBack(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return ErrQueueFull
	}
	q.items = append(q.items, item)
	return nil
}

func (q *SliceDeque) DequeueFront() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *SliceDeque) DequeueBack() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	item := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	return item, nil
}

func (q *SliceDeque) PeekFront() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return q.items[0], nil
}

func (q *SliceDeque) PeekBack() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return q.items[len(q.items)-1], nil
}

func (q *SliceDeque) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

func (q *SliceDeque) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items) == 0
}

func (q *SliceDeque) IsFull() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.capacity > 0 && len(q.items) == q.capacity
}

func (q *SliceDeque) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = make([]interface{}, 0, q.capacity)
}