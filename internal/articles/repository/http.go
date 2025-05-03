package repository

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/sashactl/internal/articles/client"
	"github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	"github.com/alexfalkowski/sashactl/internal/content"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosimple/slug"
)

const (
	noBody        = ""
	noContentType = ""
)

var bucket = aws.String("sasha-cms")

// NewRepository for books.
func NewRepository(config *config.Config, encoder *yaml.Encoder, http *client.Client, s3 *s3.Client) Repository {
	return &HTTPRepository{config: config, encoder: encoder, http: http, s3: s3}
}

// HTTPRepository uses a client to get from a site (public bucket).
type HTTPRepository struct {
	config  *config.Config
	http    *client.Client
	s3      *s3.Client
	encoder *yaml.Encoder
}

// NewArticle creates a new article with a name.
func (r *HTTPRepository) NewArticle(ctx context.Context, name string) error {
	articles, err := r.articles(ctx)
	if err != nil {
		return errors.Prefix("repository: get articles", err)
	}

	slug := slug.Make(name)

	articlesDir := filepath.Join(r.config.Path, "articles")
	articlesConfig := filepath.Join(articlesDir, "articles.yml")
	articleDir := filepath.Join(articlesDir, slug)
	articleConfig := filepath.Join(articleDir, "article.yml")

	if err := os.MkdirAll(filepath.Join(articleDir, "images"), 0o777); err != nil {
		return errors.Prefix("repository: mkdir", err)
	}

	article := &model.Article{Name: name, Slug: slug}
	articles.Articles = append(articles.Articles, article)

	configFile, err := os.Create(articlesConfig)
	if err != nil {
		return errors.Prefix("repository: create articles", err)
	}

	if err := r.encoder.Encode(configFile, articles); err != nil {
		return errors.Prefix("repository: encode articles", err)
	}

	article.Body = "Add my story!"
	article.Images = []*model.Image{
		{Name: "dummy", Description: "Add me!"},
	}

	articleFile, err := os.Create(articleConfig)
	if err != nil {
		return errors.Prefix("repository: create article", err)
	}

	if err := r.encoder.Encode(articleFile, article); err != nil {
		return errors.Prefix("repository: encode article", err)
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
		return errors.Prefix("repository: walk images", err)
	}

	return nil
}

func (r *HTTPRepository) uploadConfig(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, "articles.yml", content.YAMLContentType, path); err != nil {
		return errors.Prefix("repository: create config", err)
	}

	if err := r.put(ctx, slug+"/", noContentType, noBody); err != nil {
		return errors.Prefix("repository: create folder", err)
	}

	return nil
}

func (r *HTTPRepository) uploadArticle(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, slug+"/article.yml", content.YAMLContentType, path); err != nil {
		return errors.Prefix("repository: create article", err)
	}

	return nil
}

func (r *HTTPRepository) uploadImages(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, slug+"/images/", noContentType, noBody); err != nil {
		return errors.Prefix("repository: create images", err)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Prefix("repository: walk image", err)
		}

		if info.IsDir() {
			return nil
		}

		if err := r.put(ctx, slug+"/images/"+filepath.Base(path), content.JPEGContentType, path); err != nil {
			return errors.Prefix("repository: create images", err)
		}

		return nil
	})
}

func (r *HTTPRepository) put(ctx context.Context, path, contentType, body string) error {
	file, err := r.file(body)
	if err != nil {
		return err
	}

	defer r.close(file)

	input := &s3.PutObjectInput{
		Bucket: bucket,
		Key:    aws.String(path),
		Body:   file,
	}

	if !strings.IsEmpty(contentType) {
		input.ContentType = aws.String(contentType)
	}

	_, err = r.s3.PutObject(ctx, input)

	return err
}

func (r *HTTPRepository) file(body string) (io.ReadCloser, error) {
	if strings.IsEmpty(body) {
		return nil, nil
	}

	return os.Open(body)
}

func (r *HTTPRepository) close(rc io.ReadCloser) error {
	if rc == nil {
		return nil
	}

	return rc.Close()
}

func (r *HTTPRepository) articles(ctx context.Context) (*model.Articles, error) {
	site := &model.Articles{}

	if err := r.http.Get(ctx, r.config.Address+"/articles.yml", site); err != nil {
		if status.Code(err) == http.StatusNotFound {
			return site, nil
		}

		return nil, err
	}

	return site, nil
}
