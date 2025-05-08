package notionmd

type Config struct {
	Notion NotionConfig `yaml:"notion"`
}

type NotionConfig struct {
	Token           string `yaml:"token"`
	ParseChildPages bool   `yaml:"parse_child_pages"`
}

func NewConfig(notion *NotionConfig) (*Config, error) {
	return &Config{
		Notion: *notion,
	}, nil
}
