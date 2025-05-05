package config

import "github.com/alexfalkowski/go-service/os"

// Config for articles.
type Config struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty" toml:"path,omitempty"`
}

// GetPath from config.
func (c *Config) GetPath() string {
	return os.CleanPath(c.Path)
}
