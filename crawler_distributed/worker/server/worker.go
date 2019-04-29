package main

import (
	"flag"
	"fmt"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
	"gopcp.v2/chapter7/crawler_distributed/worker"
	"log"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Print("must speictya port")
		return
	}
	log.Fatal(rpc_support.ServeRpc(fmt.Sprintf(":%d", *port), worker.CrawlService{}))
}
