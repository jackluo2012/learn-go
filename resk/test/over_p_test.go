package test

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"testing"
)

//基准测试代码，Benchmark 开头，第一个字母B大写
//testing.B
func BenchmarkParallelUpdateForLock(b *testing.B) {
	g := GoodsSigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(1000000)
	_, err := db.Insert(g)
	if err != nil {
		fmt.Println(err)
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UpdateForLock(g.Goods)
		}
	})

}

//基准测试：无符号字段类型+直接更新
func BenchmarkParallelUpdateForUnsigned(b *testing.B) {
	g := GoodsUnsigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(1000000)
	_, err := db.Insert(g)
	if err != nil {
		fmt.Println(err)
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UpdateForUnsigned(g.Goods)
		}
	})
}

//乐观锁的基准测试
func BenchmarkParallelUpdateOptimistic(b *testing.B) {

	g := GoodsSigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(1000000)
	_, err := db.Insert(g)
	if err != nil {
		fmt.Println(err)
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UpdateForOptimistic(g.Goods)
		}
	})

}

//乐观锁+无符号字段的基准测试
func BenchmarkParallelUpdateForOptimisticAndUnsigned(b *testing.B) {
	g := GoodsUnsigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(1000000)
	_, err := db.Insert(g)
	if err != nil {
		fmt.Println(err)
		return
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UpdateForOptimisticAndUnsigned(g.Goods)
		}
	})

}
