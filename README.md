# queue
golang queue usages

# 队列使用方法示例

```go
package main

import (
	"fmt"
	"queue"
)

func main() {
	// 1. 基本队列使用
	fmt.Println("=== 基本队列 ===")
	basicQueue := queue.NewSliceQueue(5)
	_ = basicQueue.Enqueue(1)
	_ = basicQueue.Enqueue(2)
	_ = basicQueue.Enqueue(3)
	
	item, _ := basicQueue.Dequeue()
	fmt.Println("Dequeued:", item) // 1
	
	peekItem, _ := basicQueue.Peek()
	fmt.Println("Peek:", peekItem) // 2
	
	fmt.Println("Size:", basicQueue.Size()) // 2
	
	// 2. 优先级队列使用
	fmt.Println("\n=== 优先级队列 ===")
	priorityQueue := queue.NewPriorityQueue(5)
	_ = priorityQueue.Enqueue("低优先级", 3)
	_ = priorityQueue.Enqueue("高优先级", 1)
	_ = priorityQueue.Enqueue("中优先级", 2)
	
	highItem, _ := priorityQueue.Dequeue()
	fmt.Println("Dequeued:", highItem) // 高优先级
	
	midItem, _ := priorityQueue.Dequeue()
	fmt.Println("Dequeued:", midItem) // 中优先级
	
	lowItem, _ := priorityQueue.Dequeue()
	fmt.Println("Dequeued:", lowItem) // 低优先级
	
	// 3. 双端队列使用
	fmt.Println("\n=== 双端队列 ===")
	deque := queue.NewSliceDeque(5)
	_ = deque.EnqueueFront(1)
	_ = deque.EnqueueBack(2)
	_ = deque.EnqueueFront(0)
	
	frontItem, _ := deque.DequeueFront()
	fmt.Println("Dequeued front:", frontItem) // 0
	
	backItem, _ := deque.DequeueBack()
	fmt.Println("Dequeued back:", backItem) // 2
	
	// 4. 阻塞队列使用
	fmt.Println("\n=== 阻塞队列 ===")
	blockingQueue := queue.NewBlockingQueue(3)
	
	// 生产者
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Enqueueing:", i)
			blockingQueue.Enqueue(i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// 消费者
	go func() {
		for i := 0; i < 5; i++ {
			item := blockingQueue.Dequeue()
			fmt.Println("Dequeued:", item)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	
	time.Sleep(2 * time.Second)
	
	// 5. 循环队列使用
	fmt.Println("\n=== 循环队列 ===")
	circularQueue := queue.NewCircularQueue(3)
	_ = circularQueue.Enqueue("a")
	_ = circularQueue.Enqueue("b")
	_ = circularQueue.Enqueue("c")
	
	// 队列已满，无法再添加
	err := circularQueue.Enqueue("d")
	if err != nil {
		fmt.Println("Enqueue error:", err) // Queue is full
	}
	
	item1, _ := circularQueue.Dequeue()
	fmt.Println("Dequeued:", item1) // a
	
	// 现在可以添加新元素
	_ = circularQueue.Enqueue("d")
	
	item2, _ := circularQueue.Dequeue()
	fmt.Println("Dequeued:", item2) // b
}
```

# 功能总结

## 这个队列封装提供了以下功能：

### 多种队列实现：

- 基于切片的队列

- 基于链表的队列

- 优先级队列

- 循环队列

- 双端队列

- 阻塞队列

### 线程安全：所有队列实现都使用读写锁保证线程安全

### 完整操作：

- 入队(Enqueue)

- 出队(Dequeue)

- 查看队首(Peek)

- 获取队列大小(Size)

- 检查空/满状态(IsEmpty/IsFull)

- 清空队列(Clear)

- 双端队列特有操作(EnqueueFront/DequeueBack等)

- 阻塞队列特有操作(阻塞/超时入队出队)

### 错误处理：

- 队列空/满时的明确错误返回

### 扩展性：

- 通过接口设计，可以轻松添加新的队列实现

# Golang 队列类型选择指南

## 1. 队列类型概述
在编程中，队列是一种先进先出(FIFO)的数据结构，但在不同场景下，我们需要不同类型的队列来满足特定需求。以下是各种队列类型的特性及适用场景分析。

## 2. 基础队列比较

### 2.1 切片队列 (SliceQueue)
**实现特点：**
- 基于Go的slice实现
- 动态扩容(当capacity=0时)
- 内存连续

**优势：**
- 实现简单直观
- 内存局部性好，访问速度快
- 适合小规模数据

**劣势：**
- 频繁的入队出队会导致内存重新分配和复制
- 当容量固定时，出队操作会导致底层数组前部空间浪费

**适用场景：**
- 数据量不大且变化不频繁的场景
- 需要快速实现的简单队列需求
- 作为其他队列的基础实现

**示例代码：**
```go
q := queue.NewSliceQueue(0) // 无限容量
_ = q.Enqueue("item")
item, _ := q.Dequeue()
```

### 2.2 链表队列 (LinkedListQueue)
**实现特点：**
- 基于链表节点实现
- 每个元素独立分配内存

**优势：**
- 真正的O(1)时间复杂度的入队出队
- 没有容量限制(除非特别指定)
- 不会造成内存浪费

**劣势：**
- 内存不连续，缓存不友好
- 每个元素需要额外的指针存储空间
- 内存分配开销较大

**适用场景：**
- 频繁入队出队且数据量大的场景
- 无法预估队列最大容量的情况
- 需要严格O(1)时间复杂度的场景

**示例代码：**
```go
q := queue.NewLinkedListQueue(0) // 无限容量
_ = q.Enqueue("item")
item, _ := q.Dequeue()
```

## 3. 特殊队列比较

### 3.1 循环队列 (CircularQueue)
**实现特点：**
- 基于固定大小的数组实现
- 头尾指针循环移动

**优势：**
- 内存利用率高
- 入队出队都是O(1)时间复杂度
- 适合固定大小的队列需求

**劣势：**
- 容量固定，无法动态扩容
- 实现相对复杂

**适用场景：**
- 已知队列最大容量的场景
- 需要高效利用内存的场景
- 嵌入式或内存受限环境

**示例代码：**
```go
q := queue.NewCircularQueue(100) // 固定容量100
_ = q.Enqueue("item")
item, _ := q.Dequeue()
```

### 3.2 优先级队列 (PriorityQueue)
**实现特点：**
- 基于堆(通常是二叉堆)实现
- 元素按优先级出队

**优势：**
- 最高优先级的元素总是先出队
- 插入和删除都是O(log n)时间复杂度

**劣势：**
- 实现复杂度高
- 比普通队列性能稍差

**适用场景：**
- 任务调度系统
- 需要处理优先级的场景
- Dijkstra等算法实现

**示例代码：**
```go
q := queue.NewPriorityQueue(0) // 无限容量
_ = q.Enqueue("低优先级", 3)
_ = q.Enqueue("高优先级", 1)
item, _ := q.Dequeue() // 返回"高优先级"
```

### 3.3 双端队列 (Deque)
**实现特点：**
- 两端都可入队出队
- 基于切片或链表实现

**优势：**
- 操作灵活
- 可以同时作为队列和栈使用

**劣势：**
- 实现比普通队列复杂
- 某些操作性能可能受影响

**适用场景：**
- 需要两端操作的场景
- 滑动窗口算法
- 撤销操作历史记录

**示例代码：**
```go
q := queue.NewSliceDeque(0) // 无限容量
_ = q.EnqueueFront("front") // 从前面入队
_ = q.EnqueueBack("back")   // 从后面入队
front, _ := q.DequeueFront()
back, _ := q.DequeueBack()
```

### 3.4 阻塞队列 (BlockingQueue)
**实现特点：**
- 基于条件变量实现
- 线程安全的阻塞操作

**优势：**
- 完美的生产者-消费者模型
- 自动的流量控制
- 线程间通信利器

**劣势：**
- 性能比非阻塞队列差
- 可能造成死锁

**适用场景：**
- 多线程/协程环境
- 生产者-消费者模式
- 需要流量控制的场景

**示例代码：**
```go
q := queue.NewBlockingQueue(10)

// 生产者
go func() {
    for i := 0; i < 100; i++ {
        q.Enqueue(i)
    }
}()

// 消费者
go func() {
    for {
        item := q.Dequeue()
        fmt.Println(item)
    }
}()
```

## 4. 性能对比

| 队列类型 | 入队时间复杂度 | 出队时间复杂度 | 空间复杂度 | 线程安全 | 特性 |
|---------|-------------|-------------|---------|-------|-----|
| SliceQueue | O(1)* | O(n) | O(n) | 是 | 简单实现，动态扩容 |
| LinkedListQueue | O(1) | O(1) | O(n) | 是 | 真正O(1)操作 |
| CircularQueue | O(1) | O(1) | O(n) | 是 | 固定容量，高效 |
| PriorityQueue | O(log n) | O(log n) | O(n) | 是 | 按优先级出队 |
| Deque | O(1)* | O(1)* | O(n) | 是 | 双端操作 |
| BlockingQueue | O(1) | O(1) | O(n) | 是 | 阻塞操作 |

*注：SliceQueue和Deque的复杂度在动态扩容时为O(n)*

## 5. 选择建议
- **默认选择**：如果不确定使用哪种队列，SliceQueue或LinkedListQueue是不错的选择
- **固定大小需求**：如果需要固定大小的队列，选择CircularQueue
- **优先级处理**：需要按优先级处理元素时，选择PriorityQueue
- **双端操作**：需要从两端操作队列时，选择Deque
- **多线程通信**：在goroutine间传递数据时，选择BlockingQueue
- **性能敏感**：对性能要求极高且数据量大时，选择LinkedListQueue或CircularQueue
- **内存敏感**：在内存受限环境中，选择CircularQueue

## 6. 最佳实践

### 6.1 容量规划：
- 如果知道队列的最大容量，创建时指定容量可以提高性能
- 特别是对于切片实现的队列，避免频繁扩容

### 6.2 错误处理：
- 总是检查Enqueue和Dequeue的返回错误
- 特别是对于固定容量的队列

### 6.3 并发控制：
- 虽然这些实现都是线程安全的，但复杂操作还是需要额外同步
- 例如：检查大小后再操作这种复合操作

### 6.4 资源释放：
- 对于存储指针或资源的队列，出队后记得释放资源
- 特别是当队列作为对象池使用时

### 6.5 性能监控：
- 在高性能场景中，监控队列长度和操作耗时
- 根据实际情况调整队列类型或容量

## 7. 实际应用案例

### 7.1 Web服务器请求队列
```go
// 使用阻塞队列处理HTTP请求
const MaxRequests = 1000
requestQueue := queue.NewBlockingQueue(MaxRequests)

// 启动工作池
for i := 0; i < runtime.NumCPU(); i++ {
    go func() {
        for {
            req := requestQueue.Dequeue().(*http.Request)
            // 处理请求
        }
    }()
}

// 在HTTP处理器中入队
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    requestQueue.Enqueue(r)
})
```

### 7.2 任务调度系统
```go
// 使用优先级队列调度任务
taskQueue := queue.NewPriorityQueue(0)

// 添加任务
_ = taskQueue.Enqueue(Task{Name: "紧急任务"}, 1) // 高优先级
_ = taskQueue.Enqueue(Task{Name: "普通任务"}, 3) // 低优先级

// 工作线程
go func() {
    for {
        task, _ := taskQueue.Dequeue()
        // 执行任务
    }
}()
```

### 7.3 游戏事件处理
```go
// 使用双端队列处理游戏事件
eventQueue := queue.NewSliceDeque(1000)

// 添加输入事件(高优先级，插队到前面)
_ = eventQueue.EnqueueFront(PlayerInputEvent{})

// 添加普通事件(低优先级，排队到后面)
_ = eventQueue.EnqueueBack(AIEvent{})

// 游戏主循环
for {
    event, _ := eventQueue.DequeueFront()
    // 处理事件
}
```

## 8. 总结
选择正确的队列类型可以显著提高程序性能和可维护性。考虑以下因素做出选择：
- **数据量大小**：小数据量用切片队列，大数据量用链表队列
- **操作模式**：普通FIFO用基础队列，需要优先级用优先级队列
- **线程需求**：单线程用非同步队列，多线程用阻塞队列
- **内存限制**：内存敏感用循环队列，不敏感用链表队列
- **操作位置**：只需要队尾操作用普通队列，需要双端操作用双端队列