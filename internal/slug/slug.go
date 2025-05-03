package slug

import "github.com/gosimple/slug"

// Register slug settings.
func Register() {
	slug.MaxLength = 60
}
