package parser

import (
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"regexp"
)

var profileRe = regexp.MustCompile(`<tr><th><a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a></th></tr>`)
var cityUrlRe = regexp.MustCompile(`http://www.zhenai.com/zhenghun/shanghai/[^"]+`)

func ParseCity(content []byte, _ string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(content, -1)
	result := engine.ParseResult{}
	for _, match := range matches {
		url := string(match[1])
		result.Request = append(result.Request, engine.Request{
			Url:    url,
			Parser: NewProfileParse(""),
		})
	}
	//*
	matches = cityUrlRe.FindAllSubmatch(content, -1)
	for _, match := range matches {
		for _, url := range match {
			result.Request = append(result.Request, engine.Request{
				Url:    string(url),
				Parser: engine.NewFuncParser(ParseCity, config.ParseCity),
			})
		}

	}
	//*/
	return result
}
