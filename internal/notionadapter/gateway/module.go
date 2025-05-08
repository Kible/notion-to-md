package gateway

import (
	"github.com/Kible/notion-to-md/internal/notionadapter/gateway/notionapi"
	"github.com/Kible/notion-to-md/pkg/notiontomd/config"
)

type Module struct {
	NotionAPI notionapi.Gateway
}

func NewModule(cfg *config.Config) (*Module, error) {
	client := notionapi.NewClient(cfg.Notion.Token)

	notionAPI, err := notionapi.New(notionapi.Params{
		Client: client,
	})
	if err != nil {
		return nil, err
	}

	return &Module{
		NotionAPI: notionAPI,
	}, nil
}
