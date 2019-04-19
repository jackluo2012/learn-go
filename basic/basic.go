package main

import (
	"fmt"
	"math"
)

func triangle() {
	var a, b int = 3, 4
	fmt.Println(calcTringle(a, b))

}
func calcTringle(a, b int) int {
	return int(math.Sqrt(float64(a*a + b*b)))
}

func main() {
	triangle()
}
