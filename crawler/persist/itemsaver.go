package persist

import (
	"context"
	"errors"
	"github.com/olivere/elastic"
	"gopcp.v2/chapter7/crawler/engine"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {

	out := make(chan engine.Item)
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return out, err
	}
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got Item #%d, %v", itemCount, item)
			itemCount++
			if err := save(client, index, item); err != nil {
				log.Printf("Item Saver: error"+"savingitem %v:%v", item, err)
			}
		}
	}()
	return out, nil
}
func save(client *elastic.Client, index string, item engine.Item) error {
	if item.Type == "" {
		return errors.New("must supply Type")
	}
	indexserver := client.Index().Index(index).Type(item.Type)
	if item.Id != "" {
		indexserver.Id(item.Id)
	}

	if _, err := indexserver.BodyJson(item).Do(context.Background()); err != nil {
		return err
	}
	return nil
}
