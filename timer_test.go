package timewheel

import (
	"testing"
	"time"
)

func TestTimer_Stop(t *testing.T) {
	w := New(time.Second, 10)
	w.Start()

	done := false

	timer := w.AfterFunc(time.Second*2, func() {
		done = true
	})

	w.AfterFunc(time.Second, func() {
		timer.Stop()
	})

	time.Sleep(time.Second * 5)

	if done == true {
		t.Error("stop failed")
	}
}

func TestTimer_Reset(t *testing.T) {
	w := New(time.Second, 10)
	w.Start()

	var start = time.Now()
	var end time.Time
	timer := w.AfterFunc(time.Second*2, func() {
		end = time.Now()
	})

	w.AfterFunc(time.Second, func() {
		ok := timer.Reset(time.Second * 2)
		if !ok {
			t.Error("test reset failed failed")
		}
	})

	time.Sleep(time.Second * 5)

	if end.Sub(start) < time.Second*2 {
		t.Error("reset failed")
	}
}

func TestTimer_Reset_failed(t *testing.T) {
	w := New(time.Second, 10)
	w.Start()

	var start = time.Now()
	var end time.Time

	timer := w.AfterFunc(time.Second*1, func() {
		end = time.Now()
	})

	w.AfterFunc(time.Second*2, func() {
		ok := timer.Reset(time.Second * 2)

		t.Log("reset success ? :", ok)
		if ok {
			t.Error("test reset failed failed")
		}
	})

	time.Sleep(time.Second * 5)

	if end.Sub(start) < time.Second*2 {
		t.Error("reset failed")
	}
}
