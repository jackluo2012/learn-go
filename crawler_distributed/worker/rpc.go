package worker

import "gopcp.v2/chapter7/crawler/engine"

type CrawlService struct {
}

func (CrawlService) Process(req Request, result *ParseResult) error {

	engineRequet, err := DeSerializedRequest(req)
	if err != nil {
		return err
	}
	engineResult, err := engine.Worker(engineRequet)
	if err != nil {
		return err
	}
	*result = SerializedParseResult(engineResult)
	return nil
}
