package main

import (
	"flag"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/scheduler"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
	"gopcp.v2/chapter7/crawler_distributed/config"
	itemsaver "gopcp.v2/chapter7/crawler_distributed/persist/client"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
	worker "gopcp.v2/chapter7/crawler_distributed/worker/client"
	"log"
	"net/rpc"
	"strings"
)

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver hosts")
	workerHost    = flag.String("worker_hosts", "", "worker_hosts hosts")
)

func main() {
	flag.Parse()
	itemServer, err := itemsaver.ItemSaver(*itemSaverHost) // persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}

	pool := createClientPool(strings.Split(*workerHost, ","))

	//*
	processor := worker.CreateProcessor(pool)
	if err != nil {
		panic(err)
	} //*/
	e := engine.ConcurrentEngine{
		//Scheduler: &scheduler.SimpleScheduler{},//并发版的
		Scheduler:        &scheduler.QueuedScheduler{}, //队列版的
		WorkerCount:      10,
		ItemChan:         itemServer,
		RequestProcessor: processor,
		//RequestProcessor: engine.Worker,
	}

	//e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: parser.ParseCityList})
	e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun/shanghai", Parser: engine.NewFuncParser(parser.ParseCity, config.ParseCity)})
	//e.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList)})
	//e.Run(engine.Request{Url: "http://album.zhenai.com/u/1248960099", Parser: engine.NewFuncParser(parser.ParseProfile, config.ParseProfile)})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client

	for _, host := range hosts {
		log.Printf("%s",host)
		client, err := rpc_support.NewClient(host)
		if err != nil {
			log.Printf("create client err:%v", err)
		} else {
			clients = append(clients, client)
			log.Printf("create client sucess!!!")
		}
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
