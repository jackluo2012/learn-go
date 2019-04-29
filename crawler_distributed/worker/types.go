package worker

import (
	"errors"
	"fmt"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/zhenai/parser"
	"gopcp.v2/chapter7/crawler_distributed/config"
	"log"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items   []engine.Item
	Request []Request
}

func SerializedRequest(r engine.Request) Request {
	name, arg := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: arg,
		},
	}
}

//序列化
func SerializedParseResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, request := range r.Request {
		result.Request = append(result.Request, SerializedRequest(request))
	}
	return result
}

func DeSerializedRequest(r Request) (engine.Request, error) {

	parser, err := deSerializedParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

//反序列化
func DeSerializedParseResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Request {
		engineReq, err := DeSerializedRequest(req)
		if err != nil {
			log.Printf("error deserialize request %v", err)
			fmt.Print(err)
			continue
		}
		result.Request = append(result.Request, engineReq)
	}
	return result
}
func deSerializedParser(p SerializedParser) (engine.Parser, error) {

	switch p.Name {
	case config.ParseCity:
		return engine.NewFuncParser(
			parser.ParseCity,
			config.ParseCity), nil
	case config.ParseCityList:
		return engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	case config.ParseProfile:
		if name, ok := p.Args.(string); ok {
			return parser.NewProfileParse(name), nil
		} else {
			return nil, fmt.Errorf("invalid 222 args:%v", p.Args)
		}
	default:
		return nil, errors.New("unkown parser name")
	}
}
