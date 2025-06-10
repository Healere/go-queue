package queue
import (
	"sync"
)	

// LinkedListQueue 基于链表的队列实现
type LinkedListQueue struct {
	head     *node
	tail     *node
	size     int
	capacity int
	mu       sync.RWMutex
}

type node struct {
	value interface{}
	next  *node
}

// NewLinkedListQueue 创建新的链表队列
func NewLinkedListQueue(capacity int) *LinkedListQueue {
	return &LinkedListQueue{
		capacity: capacity,
	}
}

func (q *LinkedListQueue) Enqueue(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return ErrQueueFull
	}

	newNode := &node{value: item}
	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
	q.size++
	return nil
}

func (q *LinkedListQueue) Dequeue() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}

	item := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return item, nil
}

func (q *LinkedListQueue) Peek() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return q.head.value, nil
}

func (q *LinkedListQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.size
}

func (q *LinkedListQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.size == 0
}

func (q *LinkedListQueue) IsFull() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.capacity > 0 && q.size == q.capacity
}

func (q *LinkedListQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.head = nil
	q.tail = nil
	q.size = 0
}