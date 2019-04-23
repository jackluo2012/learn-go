package engine

type Request struct {
	Url        string //要请求的 url
	ParserFunc func([]byte) ParseResult //需要解析做的事情
}
type ParseResult struct {
	Request []Request //里面有一堆的 url
	Items   []interface{} //以及请求的 标签
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
