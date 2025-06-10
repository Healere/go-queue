package queue
import (
	"sync"
)
// PriorityQueue 优先级队列
type PriorityQueue struct {
	items    []*priorityItem
	capacity int
	mu       sync.RWMutex
}

type priorityItem struct {
	value    interface{}
	priority int
}

// NewPriorityQueue 创建优先级队列
func NewPriorityQueue(capacity int) *PriorityQueue {
	return &PriorityQueue{
		items:    make([]*priorityItem, 0, capacity),
		capacity: capacity,
	}
}

func (pq *PriorityQueue) Enqueue(item interface{}, priority int) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.IsFull() {
		return ErrQueueFull
	}

	newItem := &priorityItem{
		value:    item,
		priority: priority,
	}

	pq.items = append(pq.items, newItem)
	pq.heapifyUp(len(pq.items) - 1)
	return nil
}

func (pq *PriorityQueue) Dequeue() (interface{}, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.IsEmpty() {
		return nil, ErrQueueEmpty
	}

	item := pq.items[0]
	lastIndex := len(pq.items) - 1
	pq.items[0] = pq.items[lastIndex]
	pq.items = pq.items[:lastIndex]
	pq.heapifyDown(0)
	return item.value, nil
}

func (pq *PriorityQueue) Peek() (interface{}, error) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if pq.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	return pq.items[0].value, nil
}

func (pq *PriorityQueue) heapifyUp(index int) {
	for {
		parent := (index - 1) / 2
		if parent == index || pq.items[index].priority >= pq.items[parent].priority {
			break
		}
		pq.items[parent], pq.items[index] = pq.items[index], pq.items[parent]
		index = parent
	}
}

func (pq *PriorityQueue) heapifyDown(index int) {
	lastIndex := len(pq.items) - 1
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		if left <= lastIndex && pq.items[left].priority < pq.items[smallest].priority {
			smallest = left
		}
		if right <= lastIndex && pq.items[right].priority < pq.items[smallest].priority {
			smallest = right
		}
		if smallest == index {
			break
		}
		pq.items[index], pq.items[smallest] = pq.items[smallest], pq.items[index]
		index = smallest
	}
}

func (pq *PriorityQueue) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.items)
}

func (pq *PriorityQueue) IsEmpty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.items) == 0
}

func (pq *PriorityQueue) IsFull() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.capacity > 0 && len(pq.items) == pq.capacity
}

func (pq *PriorityQueue) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pq.items = make([]*priorityItem, 0, pq.capacity)
}