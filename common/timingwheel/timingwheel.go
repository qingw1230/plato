package timingwheel

import (
	"errors"
	"sync/atomic"
	"time"
	"unsafe"
)

// TimingWheel 层级时间轮的实现
type TimingWheel struct {
	tick      int64 // 时间刻度是多少 ms
	wheelSize int64 // 一层中 bucket 的数量

	interval    int64 // 当前时间轮可表示的时间范围
	currentTime int64 // 当前时间
	buckets     []*bucket
	dq          *DelayQueue // 延迟队列，只在必要时推进刻度

	// type: *TimingWheel
	// 指向下一层时间轮
	overflowWheel unsafe.Pointer

	exitC     chan struct{}
	waitGroup waitGroupWrapper
}

// NewTimingWheel 使用指定的时间刻度和大小创建层级时间轮
func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be greater than or equal to 1 ms"))
	}

	startMs := timeToMs(time.Now().UTC())

	return newTimingWheel(
		tickMs,
		wheelSize,
		startMs,
		NewDelayQueue(int(wheelSize)),
	)
}

func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, dq *DelayQueue) *TimingWheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	return &TimingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		interval:    tickMs * wheelSize,
		currentTime: truncate(startMs, tickMs),
		buckets:     buckets,
		dq:          dq,
		exitC:       make(chan struct{}),
	}
}

// add 将 t 插入到时间轮中
func (tw *TimingWheel) add(t *Timer) bool {
	curTime := atomic.LoadInt64(&tw.currentTime)
	if t.expiration < curTime+tw.tick {
		// 已经过期了
		return false
	} else if t.expiration < curTime+tw.interval {
		// 添加到当前层中
		virtualID := t.expiration / tw.tick
		b := tw.buckets[virtualID%tw.wheelSize]
		b.Add(t)

		if b.SetExpiration(virtualID * tw.tick) {
			// 该过期时间还没添加到 delayqueue，将它添加到 delayqueue 中
			tw.dq.Offer(b, b.Expiration())
		}
		return true
	} else {
		// 超出当前层可表示的范围，添加到外层时间轮
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel == nil {
			atomic.CompareAndSwapPointer(
				&tw.overflowWheel,
				nil,
				unsafe.Pointer(newTimingWheel(
					tw.interval, // 第二层的刻度是第一层表示的范围
					tw.wheelSize,
					curTime,
					tw.dq,
				)),
			)
			overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
		}
		return (*TimingWheel)(overflowWheel).add(t)
	}
}

// addOrRun 将 t 插入到时间轮中，如果已经过期则创建协程运行 t.task
func (tw *TimingWheel) addOrRun(t *Timer) {
	if !tw.add(t) {
		go t.task()
	}
}

// advanceClock 调整时间轮的当前时间
func (tw *TimingWheel) advanceClock(expiration int64) {
	curTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= curTime+tw.tick {
		curTime = truncate(expiration, tw.tick)
		atomic.StoreInt64(&tw.currentTime, curTime)

		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advanceClock(curTime)
		}
	}
}

// Start 启动时间轮
func (tw *TimingWheel) Start() {
	tw.waitGroup.Wrap(func() {
		tw.dq.Poll(tw.exitC, func() int64 {
			return timeToMs(time.Now().UTC())
		})
	})

	tw.waitGroup.Wrap(func() {
		for {
			select {
			case elem := <-tw.dq.C:
				b := elem.(*bucket)
				// 首先调整当前时间
				tw.advanceClock(b.Expiration())
				// 对时间轮的插入是相对 tw.currentTime 的
				// 二层降的一层的逻辑是：
				// 二层一格代表一层整个范围，因此 bucket 到期时，里面大部分 timer 还未到期
				// 将该 bucket 整体重新插入，已经到期的就执行，还没到期的就插入到一层了
				b.Flush(tw.addOrRun)
			case <-tw.exitC:
				return
			}
		}
	})
}

// Stop 停止当前时间轮
func (tw *TimingWheel) Stop() {
	close(tw.exitC)
	tw.waitGroup.Wait()
}

func (tw *TimingWheel) AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{
		expiration: timeToMs(time.Now().UTC().Add(d)),
		task:       f,
	}
	tw.addOrRun(t)
	return t
}

type Scheduler interface {
	Next(time.Time) time.Time
}

func (tw *TimingWheel) ScheduleFunc(s Scheduler, f func()) (t *Timer) {
	expiration := s.Next(time.Now().UTC())
	if expiration.IsZero() {
		return
	}

	t = &Timer{
		expiration: timeToMs(expiration),
		task: func() {
			expiration := s.Next(msToTime(t.expiration))
			if !expiration.IsZero() {
				t.expiration = timeToMs(expiration)
				tw.addOrRun(t)
			}

			f()
		},
	}
	tw.addOrRun(t)
	return
}
