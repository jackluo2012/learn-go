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
	p := model.Profile{
		Name:     "kimi",
		Genter:   "女士",
		Marriage: "未婚",
	}
	//fmt.Printf("%s", content)
	result := ParseProfile(content)
	p2 := result.Items[0].(model.Profile)
	if p.Name != p2.Name || p.Genter != p2.Genter || p.Marriage != p2.Marriage {
		t.Errorf("查找失败 ")
	}
	//verify result
}
