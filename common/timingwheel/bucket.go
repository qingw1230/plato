package timingwheel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Timer 定时器代表单个事件，当定时器到期时，执行指定任务
type Timer struct {
	// expiration 以 ms 为单位的过期时间
	expiration int64
	// task 过期时所执行的任务
	task func()

	// b 定时器所在的桶，type: *bucket
	b unsafe.Pointer
	// elem 所在桶内链表元素的指针
	elem *list.Element
}

// getBucket 获取所属 bucket
func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

// setBucket 设置所属 bucket
func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}

// Stop 停止该定时器
func (t *Timer) Stop() bool {
	stopped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stopped = b.Remove(t)
	}
	return stopped
}

type bucket struct {
	expiration int64

	mu     sync.Mutex
	timers *list.List
}

func newBucket() *bucket {
	return &bucket{
		expiration: -1,
		mu:         sync.Mutex{},
		timers:     list.New(),
	}
}

// Expiration 获取该 bucket 过期时间
func (b *bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

// SetExpiration 当前该 bucket 过期时间
func (b *bucket) SetExpiration(e int64) bool {
	return atomic.SwapInt64(&b.expiration, e) != e
}

// Add 添加指定定时器
func (b *bucket) Add(t *Timer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	e := b.timers.PushBack(t)
	t.setBucket(b)
	t.elem = e
}

func (b *bucket) remove(t *Timer) bool {
	if t.getBucket() != b {
		return false
	}
	b.timers.Remove(t.elem)
	t.setBucket(nil)
	t.elem = nil
	return true
}

// Remove 删除指定定时器
func (b *bucket) Remove(t *Timer) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.remove(t)
}

// Flush 对该 bucket 上的 timer 执行 reinsert 函数
func (b *bucket) Flush(reinsert func(*Timer)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for e := b.timers.Front(); e != nil; {
		next := e.Next()

		t := e.Value.(*Timer)
		b.remove(t)
		reinsert(t)

		e = next
	}

	b.SetExpiration(-1)
}
