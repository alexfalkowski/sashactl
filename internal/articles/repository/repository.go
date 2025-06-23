package repository

import "github.com/alexfalkowski/go-service/v2/context"

// Repository for articles.
type Repository interface {
	// DeleteArticle from storage.
	DeleteArticle(ctx context.Context, slug string) error

	// NewArticle to storage.
	NewArticle(ctx context.Context, name string) error

	// PublishArticle to storage.
	PublishArticle(ctx context.Context, slug string) error

	// UnpublishArticle from storage.
	UnpublishArticle(ctx context.Context, slug string) error
}
