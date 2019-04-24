package parser

import (
	"gopcp.v2/chapter7/crawler/engine"
	"regexp"
)

var profileRe = regexp.MustCompile(`<tr><th><a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a></th></tr>`)
var cityUrlRe = regexp.MustCompile(`http://www.zhenai.com/zhenghun/shanghai/[^"]+`)

func ParseCity(content []byte) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(content, -1)
	result := engine.ParseResult{}
	for _, match := range matches {
		name := string(match[2])
		result.Items = append(result.Items, "User "+name)
		result.Request = append(result.Request, engine.Request{
			Url: string(match[1]),
			ParserFunc: func(c []byte) engine.ParseResult {
				return ParseProfile(c, name)
			},
		})
	}

	matches = cityUrlRe.FindAllSubmatch(content, -1)
	for _, match := range matches {
		result.Request = append(result.Request, engine.Request{
			Url:        string(match[1]),
			ParserFunc: ParseCity,
		})
	}

	return result
}
