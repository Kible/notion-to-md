package models

type IconType string

const (
	IconEmoji    IconType = "emoji"
	IconExternal IconType = "external"
	IconFile     IconType = "file"
)

type CalloutIcon struct {
	Type     IconType          `json:"type"`
	Emoji    *string           `json:"emoji,omitempty"`
	External map[string]string `json:"external,omitempty"`
	File     map[string]string `json:"file,omitempty"`
}
