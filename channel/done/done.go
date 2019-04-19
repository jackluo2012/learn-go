package main

import (
	"fmt"
	"sync"
)

func doWorker(id int, w worker) {
	/*
	for {
		if n, ok := <-c; ok {
			fmt.Printf("Worker %d received %c\n", id, n)
		} else {
			break
		}
	}
	*/
	for n := range w.in {
		fmt.Printf("Worker %d received %c\n", id, n)
		w.done()
	}
}

/*
type worker struct {
	in   chan int
	done chan bool
}
*/

type worker struct {
	in   chan int
	done func()
}

func createWorker(id int, wg *sync.WaitGroup) worker {
	w := worker{
		in:   make(chan int, 0),
		done: func() { wg.Done() },
	}
	go doWorker(id, w)
	return w
}

func chanDemo() {
	//var c chan int //c == nil
	//c := make(chan int)
	var wg sync.WaitGroup
	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
	}
	wg.Add(20)
	//分发数据
	for i, worker := range workers {
		worker.in <- 'a' + i
	}
	//for _, worker := range workers {
	//	<-worker.done
	//}
	for i, worker := range workers {
		worker.in <- 'A' + i
	}
	wg.Wait()
	//for _, worker := range workers {
	//	<-worker.done
	//}
}

func main() {
	//
	fmt.Println("Channel as first-class citizen")
	chanDemo()

}
