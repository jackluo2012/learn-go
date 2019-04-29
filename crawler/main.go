package main

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/persist"
	"gopcp.v2/chapter7/crawler/scheduler"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
)

func main() {

	itemServer, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		//Scheduler: &scheduler.SimpleScheduler{},//并发版的
		Scheduler:        &scheduler.QueuedScheduler{}, //队列版的
		WorkerCount:      10,
		ItemChan:         itemServer,
		RequestProcessor: engine.Worker,
	}

	//e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: parser.ParseCityList})
	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun/shanghai", Parser: engine.NewFuncParser(parser.ParseCity, "ParseCity")})
}
