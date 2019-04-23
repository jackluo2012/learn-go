package engine

import (
	"log"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler // 定义调度器
	WorkerCount int       //定义处理 worker 的个数
}

//定义一个接口
type Scheduler interface {
	Submit(Request) // 向调器里 发送 Request
	ConfigureMasterWorkerChan(chan Request)  // 配送 Worker
}

func (c *ConcurrentEngine) Run(seeds ...Request) {

	in := make(chan Request) // 定义 Request in
	out := make(chan ParseResult) // 定义解析 结果
	//将 in 送入 worker chan 中
	c.Scheduler.ConfigureMasterWorkerChan(in)
	//获取 一次生成配置的 个数
	for i := 0; i < c.WorkerCount; i++ {
		createWorker(in, out)
	}
	//将请求不停的往 Submit 里面放
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
			//不停的接收 Request 的请求
			request := <-in
			//接到了就往 worker 里面放
			result, err := worker(request)
			if err != nil {
				continue
			}
			//将解析的结果 送给  out 处理
			out <- result
		}
	}()
}
