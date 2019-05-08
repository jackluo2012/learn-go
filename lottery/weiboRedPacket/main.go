package main

// curl "http://localhost:8080/set?money=1000000&uid=1&num=100000"
// curl "http://localhost:8080/get?id=3784732507&uid=1"
// wrk -t10 -c10 -d5 "http://localhost:8080/get?id=1775292212&uid=1"

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"math/rand"
	"sync"
	"time"
)

/**
 * 微博红包模拟
 */

//定一个 task 在通知中抢红包
type task struct {
	id       uint32    //id
	callback chan uint //返回
}

//红包列表
//var packageList map[uint32][]uint = make(map[uint32][]uint)
var qi = 0
var packageList *sync.Map = new(sync.Map)
var chTasks chan task = make(chan task)

const taskNum = 16

var chTasksList []chan task = make([]chan task, taskNum)

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.

	// Serve a controller based on the root Router, "/".
	mvc.New(app.Party("/")).Handle(new(lotteryController))
	return app
}

func main() {
	app := newApp()
	for i := 0; i < taskNum; i++ {

		chTasksList[i] = make(chan task, 0)

		go fetchPackageListMoney(chTasksList[i])
	}
	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) GetPing() string {
	return "pong"
}

//返回全部红包地址
// http://localhost:8080/
func (c *lotteryController) Get() map[uint32][2]int {
	rs := make(map[uint32][2]int)
	/*
	for id, list := range packageList {
		var money int
		for _, v := range list {
			money += int(v)
		}
		rs[id] = [2]int{len(list), money}

	}*/

	packageList.Range(func(key, value interface{}) bool {
		id := key.(uint32)
		list := value.([]uint)
		var money int
		for _, v := range list {
			money += int(v)
		}
		rs[id] = [2]int{len(list), money}
		return true
	})

	return rs
}

// http://localhost:8080/set?money=100&uid=1&num=100
// 发红包  设置红包
func (c *lotteryController) GetSet() string {
	uid, errUid := c.Ctx.URLParamInt("uid")           //设置红包的 id
	money, errMoney := c.Ctx.URLParamFloat64("money") //设置红包的 金额
	num, errNum := c.Ctx.URLParamInt("num")           //设置红包的  总数
	if errUid != nil || errMoney != nil || errNum != nil {
		return fmt.Sprintf("errUid=%v,errMoney=%v,errUid=%v", errUid, errMoney, errUid)
	}
	moneyTotal := int(money * 100)
	if uid < 1 || moneyTotal < num || num < 1 {
		return fmt.Sprintf("参数异常,uid=%v,errMoney=%v,errUid=%v", uid, moneyTotal, num)
	}
	//金额分配算法
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	rMax := 0.55 //随机分配的最大值
	if num > 1000 {
		rMax = 0.01
	} else if num > 100 {
		rMax = 0.1
	} else if num > 10 {
		rMax = 0.3
	}
	list := make([]uint, num)
	leftMoney := moneyTotal //剩下的金额
	leftNum := num          //剩下的数量

	//	log.Printf("money=%d,num=%d,moneyTotal=%d,leftNum=%d", money, num,moneyTotal,leftNum)
	//大循环开始,分配金额给每一个红包
	for leftNum > 0 {
		if leftNum == 1 {
			log.Printf("last=leftMoney=%d", leftMoney)
			//最后一个红包,剩余的全部金额给它
			list[num-1] = uint(leftMoney)
			break
		}
		//当钱和数量一样的时候,不能分拆
		if leftMoney == leftNum {
			for i := num - leftNum; i < num; i++ {
				list[i] = 1
			}
			break
		}
		//随机算法
		rMoney := int(float64(leftMoney-leftNum) * rMax)
		m := r.Intn(rMoney)
		if m < 1 {
			m = 1
		}
		//放入红包列表
		list[num-leftNum] = uint(m)
		leftMoney -= m
		leftNum--
		log.Printf("m=%d,leftMoney=%d,leftNum=%d", m, leftMoney, leftNum)
	}
	// 红包的唯一 ID
	id := r.Uint32()
	packageList.Store(id, list)
	//返回抢红包的 url
	return fmt.Sprintf("/get?id=%d&uid=%d&num=%d", id, uid, num)
}

//抢红包
// http://localhost:8080/get?id=3784732507&uid=1
func (c *lotteryController) GetGet() string {
	id, errUid := c.Ctx.URLParamInt("id")   //设置红包的 id
	uid, errNum := c.Ctx.URLParamInt("uid") //设置红包的  总数
	if errUid != nil || errNum != nil {
		return fmt.Sprintf("errUid=%v,errMoney=%v", errUid, errUid)
	}
	if id < 1 || uid < 1 {
		return fmt.Sprintf("红包不存在,id=%d\n", id)
	}
	//构造一个抢红包任务
	callback := make(chan uint)
	t := task{id: uint32(id), callback: callback}
	chTasks := chTasksList[qi%taskNum]
	qi++
	//发送任务
	chTasks <- t
	//返回结果
	money := <-callback
	if money <= 0 {
		return "很抱歉你没有抽中!!!"
	}
	return fmt.Sprintf("恭喜你抢到一个红包,金额为%d \n", money)
}

//启用单独的服务进行循环

func fetchPackageListMoney(chTasks chan task) {
	for {
		t := <-chTasks
		id := t.id
		l, ok := packageList.Load(id)
		if ok && l != nil {

			list := l.([]uint)
			//分配一个随机数
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			i := r.Intn(len(list))
			log.Printf("i=%d,len=%d", i, len(list))
			money := list[i]
			//更新红包中的列表信息
			if len(list) > 1 {

				if i == len(list)-1 {
					//packageList[uint32(id)] = list[:i] //过滤掉最后一个
					packageList.Store(uint32(id), list[:i]) //过滤掉最后一个
				} else if i == 0 {
					//packageList[uint32(id)] = list[1:] //过滤掉第一个
					packageList.Store(uint32(id), list[1:]) //过滤掉第一个
				} else {
					//
					packageList.Store(uint32(id), append(list[:i], list[i+1:]...))
				}
			} else {
				//如果 没有了,就删掉它
				packageList.Delete(uint32(id))
			}
			t.callback <- money
		} else {
			t.callback <- 0
		}

	}
}
