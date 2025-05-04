package config

// Config for articles.
type Config struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty" toml:"path,omitempty"`
}
