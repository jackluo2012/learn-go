package algo

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSimpRand(t *testing.T) {
	count, amount := int64(10), int64(10000)
	remain := amount
	sum := int64(0)
	for i := int64(0); i < count; i++ {
		x := SimpRand(count-i, remain)
		remain -= x
		sum += x
	}
	Convey("简单随机算法", t, func() {
		So(sum, ShouldAlmostEqual, amount)
	})

}
