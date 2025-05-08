package notionapi

import notion "github.com/jomei/notionapi"

type Client struct {
	inner *notion.Client
}

func NewClient(token string) *Client {
	return &Client{inner: notion.NewClient(notion.Token(token))}
}
