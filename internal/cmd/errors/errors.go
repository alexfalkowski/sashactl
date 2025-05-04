package errors

import "errors"

var (
	// ErrNoName when we forget to pass a name.
	ErrNoName = errors.New("new: no name provided")

	// ErrNoSlug when we forget to pass a slug.
	ErrNoSlug = errors.New("publish: no slug provided")
)
