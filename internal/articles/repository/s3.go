package repository

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/mime"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	as3 "github.com/alexfalkowski/sashactl/internal/aws/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosimple/slug"
)

var bucket = aws.String("articles")

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

	defer configFile.Close()

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
		return errors.Prefix("repository: upload config", err)
	}

	articlePath := filepath.Join(articlesPath, slug)
	articleConfig := filepath.Join(articlePath, "article.yml")

	if err := r.uploadArticle(ctx, slug, articleConfig); err != nil {
		return errors.Prefix("repository: upload article", err)
	}

	imagesPath := filepath.Join(articlePath, "images")

	if err := r.uploadImages(ctx, slug, imagesPath); err != nil {
		return errors.Prefix("repository: upload images", err)
	}

	return nil
}

// DeleteArticle from the bucket.
func (r *S3Repository) DeleteArticle(ctx context.Context, slug string) error {
	articles, err := r.articles(ctx)
	if err != nil {
		return errors.Prefix("repository: get articles", err)
	}

	articlesPath, articlesConfig := r.configPath()
	articlePath := filepath.Join(articlesPath, slug)

	if err := r.delete(ctx, articlesPath, articlePath); err != nil {
		return errors.Prefix("repository: delete files", err)
	}

	if err := os.RemoveAll(articlePath); err != nil {
		return errors.Prefix("repository: delete folder", err)
	}

	if err := r.deleteConfig(ctx, slug, articlesConfig, articles); err != nil {
		return errors.Prefix("repository: delete config", err)
	}

	return nil
}

func (r *S3Repository) uploadConfig(ctx context.Context, path string) error {
	if err := r.put(ctx, "articles.yml", mime.YAMLMediaType, path); err != nil {
		return err
	}

	return nil
}

func (r *S3Repository) uploadArticle(ctx context.Context, slug, path string) error {
	if err := r.put(ctx, filepath.Join(slug, "article.yml"), mime.YAMLMediaType, path); err != nil {
		return err
	}

	return nil
}

func (r *S3Repository) uploadImages(ctx context.Context, slug, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err := r.put(ctx, filepath.Join(slug, "images", filepath.Base(path)), mime.JPEGMediaType, path); err != nil {
			return err
		}

		return nil
	})
}

func (r *S3Repository) deleteConfig(ctx context.Context, slug, path string, articles *model.Articles) error {
	articles.Articles = slices.DeleteFunc(articles.Articles, func(a *model.Article) bool { return a.Slug == slug })

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := r.encoder.Encode(file, articles); err != nil {
		return err
	}

	if err := r.uploadConfig(ctx, path); err != nil {
		return err
	}

	return nil
}

func (r *S3Repository) configPath() (string, string) {
	articlesPath := filepath.Join(r.config.GetPath(), "articles")
	articlesConfig := filepath.Join(articlesPath, "articles.yml")

	return articlesPath, articlesConfig
}

func (r *S3Repository) delete(ctx context.Context, base, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		input := &s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    aws.String(rel),
		}

		_, err = r.s3.DeleteObject(ctx, input)

		return err
	})
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
