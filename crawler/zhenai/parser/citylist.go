package parser

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(cotent []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(cotent, -1)

	result := engine.ParseResult{}

	for _, match := range matches {
		result.Request = append(result.Request, engine.Request{
			Url:    string(match[1]),
			Parser: engine.NewFuncParser(ParseCity, config.ParseCity),
		})
	}
	return result
}
