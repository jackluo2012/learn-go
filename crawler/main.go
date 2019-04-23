package main

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
)

func main() {
	engine.Run(engine.Request{Url: "http://www.zhenai.com/zhenghun", ParserFunc: parser.ParseCityList})
}
