package config

import (
	"github.com/alexfalkowski/go-service/config"
	articles "github.com/alexfalkowski/sashactl/internal/articles/config"
)

// Config for the client.
type Config struct {
	Articles       *articles.Config `yaml:"articles,omitempty" json:"articles,omitempty" toml:"articles,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func decorateConfig(cfg *Config) *config.Config {
	return cfg.Config
}

func articlesConfig(cfg *Config) *articles.Config {
	return cfg.Articles
}
