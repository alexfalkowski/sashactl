package repository

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/strings"
	http "github.com/alexfalkowski/sashactl/internal/articles/client"
	articles "github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosimple/slug"
)

const noBody = ""

var bucket = aws.String("sasha-cms")

// NewRepository for books.
func NewRepository(config *articles.Config, encoder *yaml.Encoder, http *http.Client, s3 *s3.Client) Repository {
	return &HTTPRepository{config: config, encoder: encoder, http: http, s3: s3}
}

// HTTPRepository uses a client to get from a site (public bucket).
type HTTPRepository struct {
	config  *articles.Config
	http    *http.Client
	s3      *s3.Client
	encoder *yaml.Encoder
}

// NewArticle creates a new article with a name.
func (r *HTTPRepository) NewArticle(ctx context.Context, name string) error {
	articles, err := r.articles(ctx)
	if err != nil {
		return se.Prefix("repository: get articles", err)
	}

	slug := slug.Make(name)

	articlesDir := filepath.Join(r.config.Path, "articles")
	articlesConfig := filepath.Join(articlesDir, "articles.yml")
	articleDir := filepath.Join(articlesDir, slug)
	articleConfig := filepath.Join(articleDir, "article.yml")

	if err := os.MkdirAll(filepath.Join(articleDir, "images"), 0o777); err != nil {
		return se.Prefix("repository: mkdir", err)
	}

	article := &model.Article{Name: name, Slug: slug}
	articles.Articles = append(articles.Articles, article)

	configFile, err := os.Create(articlesConfig)
	if err != nil {
		return se.Prefix("repository: create articles", err)
	}

	if err := r.encoder.Encode(configFile, articles); err != nil {
		return se.Prefix("repository: encode articles", err)
	}

	article.Body = "Add my story!"
	article.Images = []*model.Image{
		{Name: "dummy", Description: "Add me!"},
	}

	articleFile, err := os.Create(articleConfig)
	if err != nil {
		return se.Prefix("repository: create article", err)
	}

	if err := r.encoder.Encode(articleFile, article); err != nil {
		return se.Prefix("repository: encode article", err)
	}

	return nil
}

// PublishArticle to the bucket.
func (r *HTTPRepository) PublishArticle(ctx context.Context, slug string) error {
	articlesPath := filepath.Join(r.config.Path, "articles")
	articlesConfig := filepath.Join(articlesPath, "articles.yml")

	if err := r.uploadConfig(ctx, slug, articlesConfig); err != nil {
		return err
	}

	articlePath := filepath.Join(articlesPath, slug)
	articleConfig := filepath.Join(articlePath, "article.yml")

	if err := r.uploadArticle(ctx, slug, articleConfig); err != nil {
		return err
	}

	imagesPath := filepath.Join(articlePath, "images")

	if err := r.uploadImages(ctx, slug, imagesPath); err != nil {
		return se.Prefix("repository: walk images", err)
	}

	return nil
}

func (r *HTTPRepository) uploadConfig(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, "articles.yml", path); err != nil {
		return se.Prefix("repository: create config", err)
	}

	if err := r.put(ctx, slug+"/", noBody); err != nil {
		return se.Prefix("repository: create folder", err)
	}

	return nil
}

func (r *HTTPRepository) uploadArticle(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, slug+"/article.yml", path); err != nil {
		return se.Prefix("repository: create article", err)
	}

	return nil
}

func (r *HTTPRepository) uploadImages(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, slug+"/images/", noBody); err != nil {
		return se.Prefix("repository: create images", err)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return se.Prefix("repository: walk image", err)
		}

		if info.IsDir() {
			return nil
		}

		if err := r.put(ctx, slug+"/images/"+filepath.Base(path), path); err != nil {
			return se.Prefix("repository: create images", err)
		}

		return nil
	})
}

func (r *HTTPRepository) put(ctx context.Context, path, body string) error {
	var reader io.Reader

	if !strings.IsEmpty(body) {
		file, err := os.Open(body)
		if err != nil {
			return err
		}

		defer file.Close()

		reader = file
	}

	input := &s3.PutObjectInput{Bucket: bucket, Key: aws.String(path), Body: reader}
	_, err := r.s3.PutObject(ctx, input)

	return err
}

func (r *HTTPRepository) articles(ctx context.Context) (*model.Articles, error) {
	site := &model.Articles{}

	if err := r.http.Get(ctx, r.config.Address+"/articles.yml", site); err != nil {
		return nil, err
	}

	return site, nil
}
