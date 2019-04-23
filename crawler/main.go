package main

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/scheduler"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
)

func main() {
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.SimpleScheduler{},
		WorkerCount:10,
	}

	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: parser.ParseCityList})
}
