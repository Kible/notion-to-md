package config

type ConfigInternal struct {
	Notion NotionConfig `yaml:"notion"`
}

type NotionConfig struct {
	Token           string `yaml:"token"`
	ParseChildPages bool   `yaml:"parse_child_pages"`
	ScrapeURLTitles bool   `yaml:"scrape_url_titles"`
}

func NewConfigInternal(notion *NotionConfig) (*ConfigInternal, error) {
	return &ConfigInternal{
		Notion: *notion,
	}, nil
}
