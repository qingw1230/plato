package timingwheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	// Value *bucket，过期时该 bucket 中的 timer 都过期了
	Value interface{}
	// Priority 以 expiration 作为优先级
	Priority int64
	// Index 在优先级队列中的索引
	Index int
}

// priorityQueue 由小根堆实现的优先级队列
type priorityQueue []*item

// newPriorityQueue 创建指定容量的优先级队列
func newPriorityQueue(cap int) priorityQueue {
	return make(priorityQueue, 0, cap)
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = j
	pq[j].Index = i
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(priorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}

	*pq = (*pq)[0 : n+1]
	item := x.(*item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *priorityQueue) Pop() interface{} {
	n := len(*pq)
	c := len(*pq)
	if n < (c/2) && c > 25 {
		npq := make(priorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}

	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

// PeekAndShift 获取最早过期的 item
// *item 带回过期的 item，int64 表示还有多长时间才过期
func (pq *priorityQueue) PeekAndShift(max int64) (*item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	// 检查是否过期
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}

// DelayQueue 一个无界的阻塞队列，队列头部是最早过期的 bucket
type DelayQueue struct {
	// C 保存过期的 bucket
	C chan interface{}

	mu sync.Mutex
	pq priorityQueue

	sleeping int32
	wakeupC  chan struct{}
}

func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		mu:      sync.Mutex{},
		pq:      newPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
}

// Offer 将还没添加过的 item 添加到延迟队列中
func (dq *DelayQueue) Offer(elem interface{}, expiration int64) {
	item := &item{
		Value:    elem,
		Priority: expiration,
	}

	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mu.Unlock()

	// 在优先级队列头部，说明该 bucket 最早过期，添加一个更早过期的新项
	if index == 0 {
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

// Poll 在循环中不断等待 bucket 过期，将过期 bucket 发送到 dq.C
func (dq *DelayQueue) Poll(exitC chan struct{}, nowF func() int64) {
	for {
		now := nowF()

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(now)
		if item == nil {
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			// dq 中没有任何 item 了
			if delta == 0 {
				select {
				case <-dq.wakeupC:
					// 等待新的 item 插入
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				// dq 中有 item，但都还未过期
				select {
				case <-dq.wakeupC:
					// 有更早到期的 item 插入了
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		} // if item == nil {

		// 已经有 bucket 过期了，将其添加到 dq.C 中
		select {
		case dq.C <- item.Value:
		case <-exitC:
			goto exit
		}
	} // for {

exit:
	atomic.StoreInt32(&dq.sleeping, 0)
}
