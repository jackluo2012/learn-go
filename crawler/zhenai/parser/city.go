package parser

import (
	"gopcp.v2/chapter7/crawler/engine"
	"regexp"
)

const cityRe = `<tr><th><a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a></th></tr>`

func ParseCity(cotent []byte) engine.ParseResult {
	re := regexp.MustCompile(cityRe)
	matches := re.FindAllSubmatch(cotent, -1)

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
	return result
}
