package model

// Article for site.
type Article struct {
	Name   string   `yaml:"name,omitempty"`
	Body   string   `yaml:"body,omitempty"`
	Slug   string   `yaml:"slug,omitempty"`
	Images []*Image `yaml:"images,omitempty"`
}
