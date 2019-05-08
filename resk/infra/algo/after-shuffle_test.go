package algo

import (
	"fmt"
	"testing"
)

//后洗牌算法

func TestAfterShuffle(t *testing.T) {
	count, amount := int64(10), int64(10000)
	fmt.Print(AfterShuffle(count,amount))
}
