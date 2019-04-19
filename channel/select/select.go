package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
		}

	}()
	return out
}
func worker(id int, c <-chan int) {
	for n := range c {
		fmt.Printf("Worker %d received %d\n", id, n)
		time.Sleep(time.Second)
	}
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func main() {
	//var c1, c2 chan int
	var c1, c2 = generator(), generator()
	var worker = createWorker(0)

	var values []int
	tm := time.After(time.Second * 10) //10秒后退出
	tick := time.Tick(time.Second)     //每隔1秒送一个过来
	for {
		var activeWorker chan<- int
		var activeValue int
		if len(values) > 0 {
			activeWorker = worker
			activeValue = values[0]
		}
		select {
		case n := <-c1:
			values = append(values, n)
		case n := <-c2:
			values = append(values, n)
		case activeWorker <- activeValue:
			//送完后,把每一个拿走
			values = values[1:]
		case <-time.After(800 * time.Millisecond):
			fmt.Println("太慢了吧!!!")
		case <-tm:
			fmt.Println("bye!!!")
			return
		case <-tick:
			fmt.Println("Queue Length = ", len(values))
		}
	}
}
