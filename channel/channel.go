package main

import (
	"time"
	"fmt"
)

func worker(id int, c <-chan int) {
	/*
	for {
		if n, ok := <-c; ok {
			fmt.Printf("Worker %d received %c\n", id, n)
		} else {
			break
		}
	}*/
	for n :=range c{
		fmt.Printf("Worker %d received %c\n", id, n)
	}
}
func createWorker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func chanDemo() {
	//var c chan int //c == nil
	//c := make(chan int)
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i)
	}

	//分发数据
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}
	time.Sleep(time.Millisecond)
}

func bufferChannel() {
	c := make(chan int, 3)
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'd'
}

func channelClose() {
	c := make(chan int, 3)
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'd'
	close(c)
}
func main() {
	//
	fmt.Println("Channel as first-class citizen")
	//chanDemo()
	fmt.Println("Buffered channel")

	//bufferChannel()

	fmt.Println("Channel close and range")
	channelClose()

	time.Sleep(time.Millisecond)
}
