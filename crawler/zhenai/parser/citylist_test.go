package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestParseCityList(t *testing.T) {
	content, err := ioutil.ReadFile("citylist_test_data.html")
	//content, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%s", content)
	result := ParseCityList(content)
	const resultSize = 22

	if len(result.Request) != resultSize {
		t.Errorf("result should have %d  results;but had %d", resultSize, len(result.Request))
	}

	expectedUrl := []string{
		"http://www.zhenai.com/zhenghun/aba", "http://www.zhenai.com/zhenghun/baicheng1", "http://www.zhenai.com/zhenghun/cangzhou",
	}
	expectedItems := []string{
		"City 阿坝", "City 白城", "City 沧州",
	}

	for i, url := range expectedUrl {
		if url != result.Request[i].Url {
			t.Errorf("expected url is %s,bad url is %s", url, result.Request[i].Url)
		}
	}

	for i, item := range expectedItems {
		if item != result.Items[i] {
			t.Errorf("expected url is %s,bad url is %s", item, result.Items[i])
		}
	}

	//verify result
}

func TestReaderCityList(t *testing.T) {
	resp, err := http.Get("http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s \n", content)
}
