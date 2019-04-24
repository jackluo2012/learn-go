package scheduler

import (
	"gopcp.v2/chapter7/crawler/engine"
)

type QueuedScheduler struct {
	requestChan chan engine.Request
	//每个worker 建立不同的 chan ,100个 chan 惯在
	workerChan chan chan engine.Request //每一个 worker 的chan对外 也是一个chan
}

func (q *QueuedScheduler) ConfigureMasterWorkerChan(chan engine.Request) {

}

func (q *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}
func (q *QueuedScheduler) WorkerReady(w chan engine.Request) {
	q.workerChan <- w
}

func (q *QueuedScheduler) Submit(r engine.Request) {
	q.requestChan <- r
}

/**
 * 启动 调度器
 */
func (q *QueuedScheduler) Run() {
	//初始化 通道的 通道
	q.workerChan = make(chan chan engine.Request)
	//初始化请求的通道
	q.requestChan = make(chan engine.Request)

	go func() {
		//申明 两个队列
		var requestQ []engine.Request
		var workerQ []chan engine.Request

		for {
			//生成激活的 Request 请求
			var activeRequest engine.Request
			//生成激活的 Worker
			var activeWorker chan engine.Request
			// 判断两个是否有值
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			//实现队列逻辑
			select {
			case r := <-q.requestChan:
				//send r to worker
				requestQ = append(requestQ, r)
			case w := <-q.workerChan:
				// send ? next request to w
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
