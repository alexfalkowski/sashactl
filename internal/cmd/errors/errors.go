package errors

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoName when we forget to pass a name.
	ErrNoName = errors.New("no name provided")

	// ErrNoSlug when we forget to pass a slug.
	ErrNoSlug = errors.New("no slug provided")

	// Prefix is an alias for errors.Prefix.
	Prefix = errors.Prefix
)
