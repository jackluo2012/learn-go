package engine

type ParserFunc func([]byte, string) ParseResult //需要解析做的事情

type Parser interface {
	Parse(contents []byte, url string) ParseResult
	Serialize() (name string, args interface{}) //序列化和反序列化用
}

type Request struct {
	Url    string //要请求的 url
	Parser Parser
}

//{"ParseCityList",nil} {"ParseProfile",userName}

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

//实现 nil 的
type NilParser struct {
}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}
func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

type FuncParser struct {
	parser ParserFunc //放一个函数
	name   string     //再放一个函数的名字
}

func (f *FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}

//用工厂函数来新建

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}
