package parser

import (
	"encoding/json"
	"fmt"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/model"
	"log"
	"regexp"
)

const profileJson = `<script>window.__INITIAL_STATE__.*\={"objectInfo":(.*}),"interest":`

var profileUrlRe = regexp.MustCompile(`(http://album.zhenai.com/u/[0-9]+)[">]"`)

func ParseProfile(content []byte, url string) engine.ParseResult {

	re := regexp.MustCompile(profileJson)
	matches := re.FindAllSubmatch(content, -1)
	//	log.Printf("值:%s", matches)
	result := engine.ParseResult{}
	//*
	for _, match := range matches {
		p := model.Profile{}
		err := json.Unmarshal(match[1], &p)
		if err != nil {
			log.Printf("parse Error:%s", err.Error())
		}
		//log.Printf("match 值:%v", p)
		result.Items = append(result.Items, engine.Item{
			Url:     url,
			Type:    "zhenai",
			Id:      fmt.Sprint(p.Id),
			Payload: p,
		})

	} //*/
	matches = profileUrlRe.FindAllSubmatch(content, -1)
	for _, match := range matches {
		for _, url := range match {
			result.Request = append(result.Request, engine.Request{
				Url:    string(url),
				Parser: NewProfileParse(""),
			})
		}

	}

	return result
}

type ProfileParse struct {
	userName string
}

func NewProfileParse(name string) *ProfileParse {
	return &ProfileParse{
		userName: name,
	}
}

func (p *ProfileParse) Parse(contents []byte, url string) engine.ParseResult {
	return ParseProfile(contents, url)
}

func (p *ProfileParse) Serialize() (name string, args interface{}) {
	return "ParseProfile", p.userName
}
