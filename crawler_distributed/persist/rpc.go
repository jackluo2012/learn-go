package persist

import (
	"github.com/olivere/elastic"
	"gopcp.v2/chapter7/crawler/engine"
	"gopcp.v2/chapter7/crawler/persist"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaverService) Save(item engine.Item, result *string) error {
	err := persist.Save(s.Client, s.Index, item)
	if err == nil {
		*result = "ok"
	}
	return nil
}
