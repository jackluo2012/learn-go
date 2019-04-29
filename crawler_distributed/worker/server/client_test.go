package main

import (
	"gopcp.v2/chapter7/crawler_distributed/config"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
	"gopcp.v2/chapter7/crawler_distributed/worker"
	"testing"
	"time"
)

func TestCrawlService(t *testing.T) {
	const host = ":9000"
	go rpc_support.ServeRpc(host, worker.CrawlService{})
	time.Sleep(2 * time.Second)

	client, err := rpc_support.NewClient(host)
	if err != nil {

	}
	req := worker.Request{
		Url: "http://album.zhenai.com/u/1248960099",
		Parser: worker.SerializedParser{
			Name: config.ParseProfile,
			Args: "",
		},
	}
	var result worker.ParseResult
	if client.Call(config.CrawlServiceRpc, req, &result); err != nil {
		t.Error(err)
	}
	t.Error(result)

}
