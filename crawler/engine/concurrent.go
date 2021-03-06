package engine

//先声明一个爬虫引擎
type ConcurrentEngine struct {
	Scheduler   Scheduler // 定义调度器
	WorkerCount int       //定义处理 worker 的个数

	ItemChan chan Item

	RequestProcessor Processor
}

type Processor func(Request) (ParseResult, error)

//定义一个接口
type Scheduler interface {
	ReadyNotifier
	Submit(Request)           // 向调器里 发送 Request
	WorkerChan() chan Request //问 worker chan

	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (c *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan ParseResult) // 定义解析 结果
	// 启动 调度器
	c.Scheduler.Run()
	//获取 一次生成配置的 个数
	for i := 0; i < c.WorkerCount; i++ {
		c.createWorker(c.Scheduler.WorkerChan(), out, c.Scheduler)
	}
	//将请求不停的往 Submit 里面放
	for _, r := range seeds {

		c.Scheduler.Submit(r)

	}
	for {
		//接收 parser
		result := <-out
		for _, item := range result.Items {
			// 加入存储队列
			go func() {
				c.ItemChan <- item
			}()
		}
		//再将拿到的 Request 再给调度器
		for _, r := range result.Request {
			if !isDuplicate(r.Url) {
				c.Scheduler.Submit(r)
			}
		}
	}
}

func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	//单独开个 worker 来创建
	go func() {
		for {
			//先往 里面放一个值
			ready.WorkerReady(in) //将它放入到 worker 的通道中 ,等待新的任务的到来
			//不停的接收 Request 的请求
			request := <-in
			//接到了就往 worker 里面放
			result, err := e.RequestProcessor(request) //Worker(request) //call rpc
			if err != nil {
				continue
			}
			//将解析的结果 送给  out 处理
			out <- result
		}
	}()
}

//防止重复
var visiteUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visiteUrls[url] {
		return true
	}
	visiteUrls[url] = true
	return false

}
