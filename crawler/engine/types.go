package engine

type Request struct {
	Url        string //要请求的 url
	ParserFunc ParserFunc
}

type ParserFunc func([]byte, string) ParseResult //需要解析做的事情

type ParseResult struct {
	Request []Request //里面有一堆的 url
	Items   []Item    //以及请求的 标签
}

type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
