package main

import (
	"fmt"
	"time"
)

// go run -race goroutine.go
func main() {
	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) { // race condition
			for {
				fmt.Printf("Hello from "+"goroutine %d\n", i)
			}
		}(i)
	}
	time.Sleep(time.Minute)
	fmt.Print(a)
}
