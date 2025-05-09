package models

type TextContent struct {
	Content string            `json:"content"`
	Link    map[string]string `json:"link,omitempty"`
}

type TextBlock struct {
	Type        string      `json:"type"` // always "text"
	Text        TextContent `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href,omitempty"`
}
