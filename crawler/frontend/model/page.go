package model

import "gopcp.v2/chapter7/crawler/engine"

type SearchResult struct {
	Hits  int
	Start int
	Items []engine.Item
}
