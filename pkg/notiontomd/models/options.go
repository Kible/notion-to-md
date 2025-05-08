package models

import "github.com/Kible/notion-to-md/internal/notionadapter/gateway/notionapi"

type ConfigurationOptions struct {
	SeparateChildPage     *bool `json:"separateChildPage,omitempty"`
	ConvertImagesToBase64 *bool `json:"convertImagesToBase64,omitempty"`
	ParseChildPages       *bool `json:"parseChildPages,omitempty"`
}

type NotionToMarkdownOptions struct {
	NotionClient *notionapi.Client     `json:"notionClient"`
	Config       *ConfigurationOptions `json:"config,omitempty"`
}
