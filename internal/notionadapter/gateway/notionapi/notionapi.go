package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kible/notion-to-md/internal/notionadapter/models/requests"
	"github.com/Kible/notion-to-md/internal/notionadapter/models/responses"
)

const (
	baseURL = "https://api.notion.com/v1"
)

type (
	Gateway interface {
		GetBlockChildren(context.Context, *requests.GetBlockChildrenRequest) (*responses.GetBlockChildrenResponse, error)
	}
	gateway struct {
		client *Client
	}
	Params struct {
		Client *Client
	}
)

func New(p Params) (Gateway, error) {
	return &gateway{
		client: p.Client,
	}, nil
}

func (g *gateway) GetBlockChildren(ctx context.Context, req *requests.GetBlockChildrenRequest) (*responses.GetBlockChildrenResponse, error) {
	url := fmt.Sprintf("%s/blocks/%s/children", baseURL, req.BlockID)

	if req.PageSize != nil {
		pageSize := min(*req.PageSize, 100)
		url = fmt.Sprintf("%s?page_size=%d", url, pageSize)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+string(g.client.inner.Token))
	httpReq.Header.Set("Notion-Version", "2022-06-28")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Object    string `json:"object"`
			Status    int    `json:"status"`
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("error getting block children: status %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("notion API error: status %d, code: %s, message: %s, request_id: %s",
			resp.StatusCode, errorResp.Code, errorResp.Message, errorResp.RequestID)
	}

	var notionResp *responses.GetBlockChildrenResponse
	if err := json.NewDecoder(resp.Body).Decode(&notionResp); err != nil {
		return nil, err
	}

	modifyNumberedListObject(notionResp.Results)

	return notionResp, nil
}

func modifyNumberedListObject(blocks []any) {
	numberedListIndex := 0

	for _, block := range blocks {
		if obj, ok := block.(map[string]any); ok {
			blockType, ok := obj["type"].(string)
			if !ok {
				continue
			}

			if blockType == "numbered_list_item" {
				numberedListIndex++

				if item, ok := obj["numbered_list_item"].(map[string]any); ok {
					item["number"] = numberedListIndex
				}
			} else {
				numberedListIndex = 0
			}
		}
	}
}
