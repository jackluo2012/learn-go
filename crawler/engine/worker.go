package engine

import (
	"gopcp.v2/chapter7/crawler/fetcher"
	"log"
)

func worker(r Request) (ParseResult, error) {

	log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		//如果 报错,就直接处理下个一个 url
		log.Printf("Fetcher: error fetching Url Err %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	return r.ParserFunc(body, r.Url), nil

}
