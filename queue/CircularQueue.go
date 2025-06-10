package queue
import (
	"sync"
)

// CircularQueue 循环队列
type CircularQueue struct {
	items    []interface{}
	head     int
	tail     int
	size     int
	capacity int
	mu       sync.RWMutex
}

// NewCircularQueue 创建循环队列
func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		items:    make([]interface{}, capacity),
		capacity: capacity,
	}
}

func (q *CircularQueue) Enqueue(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return ErrQueueFull
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return nil
}

func (q *CircularQueue) Dequeue() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}

	item := q.items[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return item, nil
}

func (q *CircularQueue) Peek() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return q.items[q.head], nil
}

func (q *CircularQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.size
}

func (q *CircularQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.size == 0
}

func (q *CircularQueue) IsFull() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.size == q.capacity
}

func (q *CircularQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.head = 0
	q.tail = 0
	q.size = 0
}