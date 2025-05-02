package repository

import (
	"context"

	"github.com/alexfalkowski/sashactl/internal/articles/model"
)

// Repository for books.
type Repository interface {
	// GetArticles from storage.
	GetArticles(ctx context.Context) (*model.Articles, error)
}
