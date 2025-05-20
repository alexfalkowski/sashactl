package config

import "github.com/alexfalkowski/go-service/v2/os"

// Config for articles.
type Config struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty" toml:"path,omitempty"`
}

// GetPath from config.
func (c *Config) GetPath(fs *os.FS) string {
	return fs.CleanPath(c.Path)
}
