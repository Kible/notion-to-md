package models

type EquationBlock struct {
	Type        string            `json:"type"`     // always "equation"
	Equation    map[string]string `json:"equation"` // e.g. {"expression":"E=mc^2"}
	Annotations Annotations       `json:"annotations"`
	PlainText   string            `json:"plain_text"`
	Href        interface{}       `json:"href"` // always null
}
