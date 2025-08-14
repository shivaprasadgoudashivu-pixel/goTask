package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wg *sync.WaitGroup

func main() {

	wg = new(sync.WaitGroup)
	ticker := []string{"app", "info"}

	ch1 := make(chan string, 10)
	aft := After(time.Second * 5)
	wg.Add(1)
	go sendPrice(wg, ticker, ch1, aft)
	wg.Add(1)
	go recPrice(wg, ch1, aft)
	wg.Wait()

}

func sendPrice(wg *sync.WaitGroup, ticker []string, ch1 chan string, aft chan struct{}) {
	done := false
	for i := 0; ; i++ {
		if i == len(ticker) {
			i = 0
		}
		go func() {
			<-aft
			done = true
		}()
		if done {
			close(ch1)
			wg.Done()
			runtime.Goexit()
		}
		ch1 <- ticker[i] + " : " + strconv.Itoa(rand.Intn(100))
	}
}
func recPrice(wg *sync.WaitGroup, ch1 chan string, chDuration chan struct{}) {

out:
	for {
		select {
		case price := <-ch1:
			fmt.Println(time.DateTime, price)
		case <-chDuration:
			println("Time elapses.... ")
			break out
		}
	}
	wg.Done()
}

func After(duration time.Duration) chan struct{} {
	aft := make(chan struct{})
	go func() {
		time.Sleep(duration)
		aft <- struct{}{}
		close(aft)
	}()
	return aft
}
