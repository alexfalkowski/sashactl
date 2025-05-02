package config

import (
	"github.com/alexfalkowski/go-service/config"
	articles "github.com/alexfalkowski/sashactl/internal/articles/config"
	aws "github.com/alexfalkowski/sashactl/internal/aws/config"
)

// Config for the client.
type Config struct {
	Articles       *articles.Config `yaml:"articles,omitempty" json:"articles,omitempty" toml:"articles,omitempty"`
	AWS            *aws.Config      `yaml:"aws,omitempty" json:"aws,omitempty" toml:"aws,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func decorateConfig(cfg *Config) *config.Config {
	return cfg.Config
}

func articlesConfig(cfg *Config) *articles.Config {
	return cfg.Articles
}

func awsConfig(cfg *Config) *aws.Config {
	return cfg.AWS
}
