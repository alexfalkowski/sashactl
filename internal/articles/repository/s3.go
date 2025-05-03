package repository

import (
	"context"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	as3 "github.com/alexfalkowski/sashactl/internal/aws/s3"
	"github.com/alexfalkowski/sashactl/internal/content"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosimple/slug"
)

var bucket = aws.String("sasha-cms")

// NewRepository for articles.
func NewRepository(config *config.Config, encoder *yaml.Encoder, s3 *s3.Client) Repository {
	return &S3Repository{config: config, encoder: encoder, s3: s3}
}

// S3Repository uses s3 client to interact with the content.
type S3Repository struct {
	config  *config.Config
	s3      *s3.Client
	encoder *yaml.Encoder
}

// NewArticle creates a new article with a name.
func (r *S3Repository) NewArticle(ctx context.Context, name string) error {
	articles, err := r.articles(ctx)
	if err != nil {
		return errors.Prefix("repository: get articles", err)
	}

	slug := slug.Make(name)

	articlesPath, articlesConfig := r.configPath()
	articlePath := filepath.Join(articlesPath, slug)
	articleConfig := filepath.Join(articlePath, "article.yml")

	if err := os.MkdirAll(filepath.Join(articlePath, "images"), 0o777); err != nil {
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
func (r *S3Repository) PublishArticle(ctx context.Context, slug string) error {
	articlesPath, articlesConfig := r.configPath()

	if err := r.uploadConfig(ctx, articlesConfig); err != nil {
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

func (r *S3Repository) uploadConfig(ctx context.Context, path string) error {
	if err := r.put(ctx, "articles.yml", content.YAMLContentType, path); err != nil {
		return errors.Prefix("repository: create config", err)
	}

	return nil
}

func (r *S3Repository) uploadArticle(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, slug+"/article.yml", content.YAMLContentType, path); err != nil {
		return errors.Prefix("repository: create article", err)
	}

	return nil
}

func (r *S3Repository) uploadImages(ctx context.Context, slug, path string) error {
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

func (r *S3Repository) configPath() (string, string) {
	articlesPath := filepath.Join(r.config.Path, "articles")
	articlesConfig := filepath.Join(articlesPath, "articles.yml")

	return articlesPath, articlesConfig
}

func (r *S3Repository) put(ctx context.Context, path, contentType, body string) error {
	file, err := os.Open(body)
	if err != nil {
		return err
	}

	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket:      bucket,
		Key:         aws.String(path),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	_, err = r.s3.PutObject(ctx, input)

	return err
}

func (r *S3Repository) articles(ctx context.Context) (*model.Articles, error) {
	site := &model.Articles{}
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    aws.String("articles.yml"),
	}

	out, err := r.s3.GetObject(ctx, input)
	if err != nil {
		if as3.IsNotFound(err) {
			return site, nil
		}

		return nil, err
	}

	if err := r.encoder.Decode(out.Body, site); err != nil {
		return nil, err
	}

	return site, nil
}
