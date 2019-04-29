package main

import (
	"flag"
	"fmt"
	"github.com/olivere/elastic"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"gopcp.v2/chapter7/crawler_distributed/persist"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Print("must speictya port")
		return
	}
	if err := ServeRpc(fmt.Sprintf(":%d", *port), config.ElasticIndex); err != nil {
		panic(err)
	}
}

func ServeRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	return rpc_support.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
