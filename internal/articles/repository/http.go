package repository

import (
	"context"

	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/sashactl/internal/articles/client"
	"github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
)

// NewRepository for books.
func NewRepository(config *config.Config, client *client.Client) Repository {
	return &HTTPRepository{config: config, client: client}
}

// HTTPRepository uses a client to get from a site (public bucket).
type HTTPRepository struct {
	config *config.Config
	client *client.Client
}

// GetArticles from the public bucket.
func (r *HTTPRepository) GetArticles(ctx context.Context) (*model.Articles, error) {
	site := &model.Articles{}

	if err := r.client.Get(ctx, r.config.Address+"/articles.yml", site); err != nil {
		return nil, se.Prefix("repository: get articles", err)
	}

	return site, nil
}
