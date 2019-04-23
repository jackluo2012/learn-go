package engine

import (
	"log"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler //定义调度器
	WorkerCount int       //定义好多个 worker 处理
}

//定义一个接口
type Scheduler interface {
	Submit(Request) //发送
	ConfigureMasterWorkerChan(chan Request)
}

func (c *ConcurrentEngine) Run(seeds ...Request) {

	in := make(chan Request)
	out := make(chan ParseResult)
	//将 in 送入 worker chan 中
	c.Scheduler.ConfigureMasterWorkerChan(in)
	for i := 0; i < c.WorkerCount; i++ {
		createWorker(in, out)
	}

	for _, r := range seeds {
		c.Scheduler.Submit(r)
	}
	itemCount := 0
	for {
		//接收 parser
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item $%d: %v", itemCount, item)
			itemCount++
		}
		//再将拿到的 Request 再给调度器
		for _, r := range result.Request {
			c.Scheduler.Submit(r)
		}
	}
}

func createWorker(in chan Request, out chan ParseResult) {
	//单独开个 worker 来创建
	go func() {
		for {
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			//将解析的结果 送给  out 处理
			out <- result
		}
	}()
}