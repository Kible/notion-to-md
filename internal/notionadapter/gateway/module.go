package gateway

import (
	"github.com/Kible/notion-to-md/internal/notionadapter/config"
	"github.com/Kible/notion-to-md/internal/notionadapter/gateway/notionapi"
)

type Module struct {
	NotionAPI notionapi.Gateway
}

func NewModule(cfg *config.ConfigInternal) (*Module, error) {
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
