package parser

import (
	"gopcp.v2/chapter7/crawler/model"
	"io/ioutil"
	"testing"
)

func TestParseProfile(t *testing.T) {
	content, err := ioutil.ReadFile("profile_test_data.html")
	//content, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}
	name := "kimi"
	p := model.Profile{
		Name:      "kimi",
		Genter:    "女士",
		Marriage:  "未婚",
		Income:    "5001-8000元",
		Height:    "163cm",
		Age:       29,
		Education: "大学本科",
		AvatarURL: "https://photo.zastatic.com/images/photo/312241/1248960099/20663487973330052.jpg",
		Hokou:     "上海",
	}
	//fmt.Printf("%s", content)
	result := ParseProfile(content, name)
	p2 := result.Items[0].(model.Profile)
	if p.Name != p2.Name || p.Genter != p2.Genter || p.Marriage != p2.Marriage {
		t.Errorf("查找失败 ")
	}
	//verify result
}
