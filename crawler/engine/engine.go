package engine

import (
	"gopcp.v2/chapter7/crawler/fetcher"
	"log"
)

type SimpleEngine struct {
}

func (e *SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, req := range seeds {
		requests = append(requests, req)
	}
	//把第一个拿出来开始做
	for len(requests) > 0 {
		//随机取一个url值
		r := requests[0]
		requests = requests[1:]
		//改成封装好的
		parseResult, err := worker(r)
		if err != nil {
			//如果 报错,就直接处理下个一个 url
			log.Printf("Fetcher: error fetching Url Err %s: %v", r.Url, err)
			continue
		}
		//将 []Request 用...展开放入 requests 中
		requests = append(requests, parseResult.Request...)
		//打印 items

		for _, item := range parseResult.Items {
			log.Printf("Got item %v \n", item)
		}
		//将 []Request 用...展开放入 requests 中
		requests = append(requests, parseResult.Request...)
		//打印 items
	}
}

func worker(r Request) (ParseResult, error) {

	log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		//如果 报错,就直接处理下个一个 url
		log.Printf("Fetcher: error fetching Url Err %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	return r.ParserFunc(body), nil

}
