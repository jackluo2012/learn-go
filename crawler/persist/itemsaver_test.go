package persist

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/model"
	"testing"
)

func TestItemSaver(t *testing.T) {

}
func TestSave(t *testing.T) {
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
	const index = "dating_test"

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		t.Error(err)
	}
	if err := Save(client, index, expected); err != nil {
		t.Error(err)
	}

	resp, err := client.Get().Index(index).Type(expected.Type).Id(expected.Id).Do(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", resp.Source)

	var actual engine.Item

	err = json.Unmarshal([]byte(resp.Source), &actual)
	//手动转一遍
	actualProfle, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfle

	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("got %v ,expected %v", actual, expected)
	}
}
