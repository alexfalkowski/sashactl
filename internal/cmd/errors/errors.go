package errors

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoName when we forget to pass a name.
	ErrNoName = errors.New("new: no name provided")

	// ErrNoSlug when we forget to pass a slug.
	ErrNoSlug = errors.New("publish: no slug provided")

	// Prefix is just an alias for errors.Prefix.
	Prefix = errors.Prefix
)
