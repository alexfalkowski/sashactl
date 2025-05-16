package config

import (
	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/os"
)

// Config for aws.
type Config struct {
	*client.Config  `yaml:",inline" json:",inline" toml:",inline"`
	Region          string `yaml:"region,omitempty" json:"region,omitempty" toml:"region,omitempty"`
	AccessKeyID     string `yaml:"accessKeyID,omitempty" json:"accessKeyID,omitempty" toml:"accessKeyID,omitempty"`
	AccessKeySecret string `yaml:"accessKeySecret,omitempty" json:"accessKeySecret,omitempty" toml:"accessKeySecret,omitempty"`
}

// GetAccessKeyID for aws.
func (c *Config) GetAccessKeyID(fs *os.FS) ([]byte, error) {
	return fs.ReadFile(c.AccessKeyID)
}

// GetAccessKeySecret for aws.
func (c *Config) GetAccessKeySecret(fs *os.FS) ([]byte, error) {
	return fs.ReadFile(c.AccessKeySecret)
}
