package client

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
	"log"
)

func ItemSaver(host string) (chan engine.Item, error) {

	client, err := rpc_support.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got Item #%d, %v", itemCount, item)
			itemCount++
			var result string
			//	Call RPC to Save item
			err = client.Call(config.ItemSaverRpc, item, &result)
			if err != nil && result != "ok" {
				log.Printf("Item Saver: error"+"savingitem %v:%v", item, err)
			}

		}
	}()
	return out, nil
}
