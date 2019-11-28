package timewheel

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	timer := time.NewTicker(time.Second)

	flag := false
	for {
		select {
		case <-timer.C:
			fmt.Println(time.Now())
			if !flag {
				time.Sleep(time.Second * 3)
				flag = true
			}
		default:
		}
	}
}

func Test2(t *testing.T) {
	w := NewWheel(time.Second, 5)

	log.Println("hello")

	for i := 0; i < 6; i++ {
		i := i
		w.AfterFunc(time.Second*time.Duration(i), func() {
			log.Println("hello world: ", i)
		})
		w.AfterFunc(time.Second*time.Duration(i), func() {
			log.Println("hello world: ", i)
		})
	}

	w.Start()

	//time.Sleep(time.Second *12)
	w.CloseAndWait()
}

func TestWheel_NewTimer(t *testing.T) {
	w := NewWheel(time.Second, 5)
	w.Start()

	log.Println()
	timer := w.NewTimer(time.Second)

	<-timer.C

	log.Println()
}

func TestRaceDetect(t *testing.T) {
	w := NewWheel(time.Second, 5)

	for i := 0; i < 40000; i++ {
		i := i
		w.AfterFunc(time.Duration(i/10000)*time.Second, func() {
			log.Println(i/10000)
		})
	}
	w.Start()
	w.CloseAndWait()

	time.Sleep(time.Second)
}
