package main

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/model"
	"gopcp.v2/chapter7/crawler_distributed/rpc_support"
	"testing"
	"time"
)

func TestItemSaver(t *testing.T) {
	//Start ItemServerServer
	const host = ":1234"

	go ServeRpc(host, "dating_test")
	time.Sleep(2*time.Second)
	//Start ItemServerClient
	client, err := rpc_support.NewClient(host)
	if err != nil {
		t.Errorf("%v", err)
	}
	var result string
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/1248960099",
		Id:   "1248960099",
		Type: "zhenai",
		Payload: model.Profile{
			Name:      "kimi",
			Genter:    "女士",
			Marriage:  "未婚",
			Income:    "5001-8000元",
			Height:    "163cm",
			Age:       29,
			Education: "大学本科",
			AvatarURL: "https://photo.zastatic.com/images/photo/312241/1248960099/20663487973330052.jpg",
			Hokou:     "上海",
		},
	}
	//Call save
	err = client.Call("ItemSaverService.Save", expected, &result)
	if err != nil && result != "ok" {
		t.Error(err)
	}

}
