package client

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"gopcp.v2/chapter7/crawler_distributed/worker"
	"net/rpc"
)

func CreateProcessor(clientchChan chan *rpc.Client) engine.Processor {

	return func(r engine.Request) (engine.ParseResult, error) {
		//先进行序列化
		sReq := worker.SerializedRequest(r)

		var sResult worker.ParseResult
		c := <-clientchChan
		err := c.Call(config.CrawlServiceRpc, sReq, &sResult)

		if err != nil {
			return engine.ParseResult{}, err
		} else {
			return worker.DeSerializedParseResult(sResult), nil
		}

	}
}
