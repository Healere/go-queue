package queue

import (
	"sync"
	"time"
)

// BlockingQueue 阻塞队列
type BlockingQueue struct {
	queue    Queue
	capacity int
	enqueueCond *sync.Cond
	dequeueCond *sync.Cond
	mu       sync.Mutex
}

// NewBlockingQueue 创建阻塞队列
func NewBlockingQueue(capacity int) *BlockingQueue {
	bq := &BlockingQueue{
		queue:    NewSliceQueue(capacity),
		capacity: capacity,
	}
	bq.enqueueCond = sync.NewCond(&bq.mu)
	bq.dequeueCond = sync.NewCond(&bq.mu)
	return bq
}

func (bq *BlockingQueue) Enqueue(item interface{}) {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	for bq.queue.IsFull() {
		bq.enqueueCond.Wait()
	}

	_ = bq.queue.Enqueue(item)
	bq.dequeueCond.Signal()
}

func (bq *BlockingQueue) Dequeue() interface{} {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	for bq.queue.IsEmpty() {
		bq.dequeueCond.Wait()
	}

	item, _ := bq.queue.Dequeue()
	bq.enqueueCond.Signal()
	return item
}

func (bq *BlockingQueue) TryEnqueue(item interface{}, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	bq.mu.Lock()
	defer bq.mu.Unlock()

	for bq.queue.IsFull() {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return false
		}
		// bq.enqueueCond.WaitTimeout(remaining)
	}

	_ = bq.queue.Enqueue(item)
	bq.dequeueCond.Signal()
	return true
}

func (bq *BlockingQueue) TryDequeue(timeout time.Duration) (interface{}, bool) {
	deadline := time.Now().Add(timeout)
	bq.mu.Lock()
	defer bq.mu.Unlock()

	for bq.queue.IsEmpty() {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, false
		}
		// bq.dequeueCond.WaitTimeout(remaining)
	}

	item, _ := bq.queue.Dequeue()
	bq.enqueueCond.Signal()
	return item, true
}

func (bq *BlockingQueue) Size() int {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	return bq.queue.Size()
}

func (bq *BlockingQueue) IsEmpty() bool {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	return bq.queue.IsEmpty()
}

func (bq *BlockingQueue) IsFull() bool {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	return bq.queue.IsFull()
}

func (bq *BlockingQueue) Clear() {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	bq.queue.Clear()
	bq.enqueueCond.Broadcast()
}