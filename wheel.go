package timewheel

import (
	"github.com/kuhufu/timewheel/pq"
	"log"
	"sync"
	"time"
)

/**
时间轮的相关资料
http://www.10tiao.com/html/249/201703/2651959961/1.html
*/

type Wheel struct {
	//每一刻度的时间长度
	interval time.Duration

	ticker *time.Ticker

	//时间轮槽数量
	wheelSize int64

	//环形队列
	slots []*pq.MinPQ

	//当前所在时间轮的位置索引
	curIdx int64

	//当前时间轮是第几轮
	curCycleNum int64

	//通知关闭时间轮
	done chan struct{}

	//只运行一次
	sync.Once

	//是否已关闭
	closed bool

	wg sync.WaitGroup

	mu sync.Mutex

	//运行开始时间，用来计算优先级
	createAt time.Time
}

var comparator = func(a, b interface{}) int {
	pa := a.(*pq.Item).Key.(int64)
	pb := b.(*pq.Item).Key.(int64)

	switch {
	case pa > pb:
		return 1
	case pa < pb:
		return -1
	}
	return 0
}

type Task struct {
	cycleNum int64
	f        func()
	done     bool
}

func New(interval time.Duration, wheelSize int64) *Wheel {
	w := &Wheel{
		interval:    interval,
		wheelSize:   wheelSize,
		curCycleNum: 0,
		createAt:    time.Now(),
	}

	w.slots = make([]*pq.MinPQ, wheelSize)
	for i := int64(0); i < wheelSize; i++ {
		w.slots[i] = pq.NewMinPQ(8, comparator)
	}
	return w
}

func (w *Wheel) NewTimer(d time.Duration) *Timer {
	return w.AfterFunc(d, func() {})
}

func (w *Wheel) Start() {
	go w.Run()
}

func (w *Wheel) Run() {
	w.Do(func() {
		w.ticker = time.NewTicker(w.interval)
		w.start()
	})
}

//Wait 等待时间轮完成所有定时任务
func (w *Wheel) Wait() {
	w.wg.Wait()
}

//Close 关闭时间轮，返回值表示是否关闭成功
func (w *Wheel) Close() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.closed {
		w.closed = true
		return true
	}
	return false
}

//CloseAndWait 关闭时间轮，并等待时间轮完成所有定时任务，返回值表示是否关闭成功
func (w *Wheel) CloseAndWait() bool {
	success := w.Close()
	w.Wait()
	return success
}

func (w *Wheel) start() {
	for {
		select {
		case <-w.done:
			goto done
		case <-w.ticker.C:
			idx := w.curIdx % w.wheelSize

			for {
				tt, ok := w.slots[idx].Min()
				var t *Timer
				if ok {
					t = tt.Val.(*Timer)
					if t.task.cycleNum > w.curCycleNum {
						break
					}
					w.slots[idx].DelMin() //！！！从优先级队列中删除
					go w.doTask(t)
				} else {
					break
				}
			}

			w.curIdx++
			if w.curIdx%w.wheelSize == 0 {
				w.curCycleNum++
			}
		}
	}

done:
	log.Println("退出时间轮")
}

func (w *Wheel) doTask(t *Timer) {
	t.doTask()
	w.wg.Done()
}

//AfterFunc duration等于0时将会立刻执行
func (w *Wheel) AfterFunc(d time.Duration, f func()) *Timer {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		panic("wheel is closed")
	}
	w.mu.Unlock()

	w.wg.Add(1)

	t := &Timer{
		C:          make(chan time.Time, 1),
		ownerWheel: w,
		task: &Task{
			cycleNum: w.cycleNum(d),
			f:        f,
		},
	}

	w.getSlot(d).Insert(&pq.Item{
		Key: w.priority(d),
		Val: t,
	})
	return t
}

//key 计算任务优先级
func (w *Wheel) priority(d time.Duration) int64 {
	return time.Now().UnixNano() + int64(d)
}

func (w *Wheel) getSlot(d time.Duration) *pq.MinPQ {
	return w.slots[w.slotIdx(d)]
}

func (w *Wheel) slotIdx(duration time.Duration) int64 {
	div := duration / w.interval
	mod := duration % w.interval
	if mod > 0 {
		div++
	}

	return int64(div) % w.wheelSize
}

func (w *Wheel) cycleNum(duration time.Duration) int64 {
	div := duration / (w.interval * time.Duration(w.wheelSize))
	return int64(div) + w.curCycleNum
}
