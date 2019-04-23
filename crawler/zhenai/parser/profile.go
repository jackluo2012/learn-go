package parser

import (
	"encoding/json"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/model"
	"log"
	"regexp"
)

const profileJson = `<script>window.__INITIAL_STATE__.*\={"objectInfo":(.*}),"interest":`

func ParseProfile(cotent []byte, name string) engine.ParseResult {

	re := regexp.MustCompile(profileJson)
	matches := re.FindAllSubmatch(cotent, -1)
	//	log.Printf("值:%s", matches)
	result := engine.ParseResult{}
	//*
	for _, match := range matches {
		//log.Printf("值:%s", matches[1])
		log.Printf("match 值:%s", match[1])

		p := model.Profile{}
		err := json.Unmarshal(match[1], &p)
		if err != nil {
			log.Printf("parse Error:%s", err.Error())
		}
		p.Name = name
		log.Printf("match 值:%v", p)
		result.Items = append(result.Items, p)

	} //*/
	return result
}
