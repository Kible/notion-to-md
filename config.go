package notiontomd

type Config struct {
	Notion *NotionConfig `yaml:"notion"`
}

type NotionConfig struct {
	Token           string `yaml:"token"`             // Notion API token
	ParseChildPages bool   `yaml:"parse_child_pages"` // Parse child pages
	ScrapeURLTitles bool   `yaml:"scrape_url_titles"` // Scrape the <title> tag from the html of any bookmark, embed, link_preview, or link_to_page block types urls
}

func NewConfig(notion *NotionConfig) (*Config, error) {
	return &Config{
		Notion: notion,
	}, nil
}
