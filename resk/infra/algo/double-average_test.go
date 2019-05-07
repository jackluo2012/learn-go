package algo

import (
	"fmt"
	"testing"
)

func TestDoubleRandom(t *testing.T) {
	count, amount := int64(10), int64(10000)
	for i := int64(0); i < count; i++ {
		x := DoubleRandom(count, amount)
		fmt.Print(float64(x)/float64(100), ",")
	}
	fmt.Println()
}
