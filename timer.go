package timewheel

import (
	"sync"
	"time"
)

type Timer struct {
	C          chan time.Time
	task       *Task
	ownerWheel *Wheel
	mu         sync.Mutex
}

func (t *Timer) doTask() {
	t.mu.Lock()

	if t.task.done {
		t.mu.Unlock()
		return
	}

	t.task.done = true
	t.mu.Unlock()

	t.task.f()

	select {
	case t.C <- time.Now(): //非阻塞发送时间
	default:
	}
}

func (t *Timer) Stop() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.task.done {
		return false
	}
	t.task.done = true //timer将由wheel惰性删除
	return true
}

func (t *Timer) Reset(d time.Duration) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.task.done {
		return false
	}

	t.task.done = true //旧timer将由wheel惰性删除
	t.ownerWheel.AfterFunc(d, t.task.f)
	return true
}
