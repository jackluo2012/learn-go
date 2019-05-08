package main

import (
	"fmt"
	"gopcp.v2/chapter7/resk/infra/algo"
)

func main() {
	count, amount := int64(10), int64(10000)
	for i := int64(0); i < count; i++ {
		x := algo.SimpRand(count, amount)
		fmt.Print(float64(x)/float64(100), ",")
	}
	fmt.Println()
}
