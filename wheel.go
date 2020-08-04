package timewheel

import (
	"context"
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
	slots []*MinPQ

	//当前所在时间轮的位置索引
	curIdx int64

	//当前时间轮是第几轮
	curCycleNum int64

	//只运行一次
	sync.Once

	//是否已关闭
	closed bool

	wg sync.WaitGroup

	mu sync.Mutex

	//是否只在一个协程中执行所有任务
	onlyOneGoRoutine bool

	//运行开始时间，用来计算优先级
	createAt time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

type Task struct {
	cycleNum int64
	f        func()
	done     bool
}

func New(interval time.Duration, wheelSize int64) *Wheel {
	ctx, cancel := context.WithCancel(context.Background())
	w := &Wheel{
		interval:    interval,
		wheelSize:   wheelSize,
		curCycleNum: 0,
		createAt:    time.Now(),
		ctx:         ctx,
		cancel:      cancel,
	}

	w.slots = make([]*MinPQ, wheelSize)
	for i := int64(0); i < wheelSize; i++ {
		w.slots[i] = NewMinPQ(8)
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
func (w *Wheel) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.closed {
		w.closed = true
		w.cancel()
		return nil
	}
	return nil
}

//CloseAndWait 关闭时间轮，并等待时间轮完成所有定时任务，返回值表示是否关闭成功
func (w *Wheel) CloseAndWait() error {
	err := w.Close()
	w.Wait()
	return err
}

func (w *Wheel) start() {
	for {
		select {
		case <-w.ctx.Done():
			goto done
		case <-w.ticker.C:
			idx := w.curIdx % w.wheelSize
			slot := w.slots[idx]
			for {
				item, ok := slot.PullMinIf(func(item *Item) bool {
					return item.Val.task.cycleNum <= w.curCycleNum
				})

				//当前槽位没有当前cycle的任务
				if !ok {
					break
				}

				timer := item.Val

				//检测到小于当前周期的任务
				if timer.task.cycleNum < w.curCycleNum {
					continue
				}

				if w.onlyOneGoRoutine {
					w.doTask(timer)
				} else {
					go w.doTask(timer)
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

	w.getSlot(d).Insert(&Item{
		Key: w.priority(d),
		Val: t,
	})
	return t
}

//key 计算任务优先级
func (w *Wheel) priority(d time.Duration) int64 {
	return time.Now().UnixNano() + int64(d)
}

func (w *Wheel) getSlot(d time.Duration) *MinPQ {
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
