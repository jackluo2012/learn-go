package main

import (
	"fmt"
	"regexp"
)

const text = `my Email is net.webjoy@gmail.com
		my Email is net.webjoy@qq.com
my Email is net.webjoy@sina.com.cn


`

func main() {
	//reg := regexp.MustCompile("net.webjoy@gmail.com")
	// . + 一个或多个
	//* 是多个
	reg := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9.]+(\.[a-zA-Z0-9]+))`)

	//match := reg.FindAllString(text,-1)
	match := reg.FindAllStringSubmatch(text, -1)

	for _, sub := range match {
		fmt.Println(sub)
	}

	fmt.Print(match)
}
