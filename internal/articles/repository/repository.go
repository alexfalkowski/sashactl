package repository

import "context"

// Repository for articles.
type Repository interface {
	// NewArticle to storage.
	NewArticle(ctx context.Context, name string) error

	// PublishArticle to storage.
	PublishArticle(ctx context.Context, slug string) error
}
