package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	response, err := http.Get("http://www.zhenai.com/zhenghun/chengdu/nv")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: status code	", response.StatusCode)
	}
	//将网页的 gbk 编号 换成 utf-8编码
	//utf8Reader := transform.NewReader(response.Body,simplifiedchinese.GBK.NewDecoder())

	e := determineEncoding(response.Body)
	utf8Reader := transform.NewReader(response.Body, e.NewDecoder())
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", all)
}

/**
 * 设置1024
 * 自动转换成 对应的字符编码
 */
func determineEncoding(r io.Reader) encoding.Encoding {
	//默认只读1024个字节,读了以后就不能再读
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
