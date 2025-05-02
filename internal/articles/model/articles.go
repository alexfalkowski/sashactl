package model

// Articles for site.
type Articles struct {
	Articles []*Article `yaml:"articles,omitempty"`
}
