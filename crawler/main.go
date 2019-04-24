package main

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
)

func main() {
	e := engine.ConcurrentEngine{
		//Scheduler: &scheduler.SimpleScheduler{},//并发版的
		//Scheduler: &scheduler.QueuedScheduler{},//队列版的
		WorkerCount:10,
	}

	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: parser.ParseCityList})
}
