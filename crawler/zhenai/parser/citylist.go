package parser

import (
	"gopcp.v2/chapter7/crawler/engine"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(cotent []byte) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(cotent, -1)

	result := engine.ParseResult{}

	for _, match := range matches {
		result.Items = append(result.Items, "City "+string(match[2]))
		result.Request = append(result.Request, engine.Request{
			Url:        string(match[1]),
			ParserFunc: ParseCity,
		})
	}
	return result
}
