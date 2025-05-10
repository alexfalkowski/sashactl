package model

// Article for site.
type Article struct {
	Name string `yaml:"name,omitempty"`
	Slug string `yaml:"slug,omitempty"`
}
