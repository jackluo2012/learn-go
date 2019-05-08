package algo

import (
	"fmt"
	"testing"
)

func TestDoubleAverage(t *testing.T) {
	count, amount := int64(10), int64(10000)
	for i := int64(0); i < count; i++ {
		x := DoubleAverage(count, amount)
		fmt.Print(float64(x)/float64(100), ",")
	}
	fmt.Println()
}
