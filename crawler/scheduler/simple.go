package scheduler

import "gopcp.v2/chapter7/crawler/engine"

type SimpleScheduler struct {
	workerChan chan engine.Request //调度里有一个通道
}

//设置初始化的 workerchan
func (s *SimpleScheduler) WorkerReady(c chan engine.Request)  {

}
func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

// 设置往 workerChan 写的 reques
func (s *SimpleScheduler) Submit(r engine.Request) {
	// send request down to worker chan
	go func() {
		s.workerChan <- r
	}()

}
