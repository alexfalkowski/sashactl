package config

import "github.com/alexfalkowski/go-service/client"

// Config for articles.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
	Path           string `yaml:"path,omitempty" json:"path,omitempty" toml:"path,omitempty"`
}
