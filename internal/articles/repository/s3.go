package repository

import (
	"io"
	"slices"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
	"github.com/alexfalkowski/sashactl/internal/articles/config"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	"github.com/alexfalkowski/sashactl/internal/aws/s3"
	"github.com/gosimple/slug"
)

var bucket = ptr.Value("articles")

// Params for articles.
type Params struct {
	di.In

	Config    *config.Config
	Encoder   *yaml.Encoder
	S3        *s3.Client
	FS        *os.FS
	Generator id.Generator
}

// NewRepository for articles.
func NewRepository(params Params) Repository {
	return &S3Repository{
		config:    params.Config,
		encoder:   params.Encoder,
		s3:        params.S3,
		fs:        params.FS,
		generator: params.Generator,
	}
}

// S3Repository uses s3 client to interact with the content.
type S3Repository struct {
	config    *config.Config
	s3        *s3.Client
	encoder   *yaml.Encoder
	fs        *os.FS
	generator id.Generator
}

// DeleteArticle from disk.
func (r *S3Repository) DeleteArticle(ctx context.Context, slug string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("repository: delete article", runtime.ConvertRecover(r))
		}
	}()

	ctx = meta.WithRequestID(ctx, meta.String(r.generator.Generate()))
	articles := r.articles(ctx)
	articlesPath, articlesConfig := r.configPath()
	articlePath := r.fs.Join(articlesPath, slug)

	err = r.fs.RemoveAll(articlePath)
	runtime.Must(err)

	r.deleteConfig(ctx, slug, articlesConfig, articles)

	return nil
}

// NewArticle creates a new article with a name.
func (r *S3Repository) NewArticle(ctx context.Context, name string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("repository: new article", runtime.ConvertRecover(r))
		}
	}()

	ctx = meta.WithRequestID(ctx, meta.String(r.generator.Generate()))
	articles := r.articles(ctx)
	slug := slug.Make(name)
	articlesPath, articlesConfig := r.configPath()
	articlePath := r.fs.Join(articlesPath, slug)
	articleConfigPath := r.fs.Join(articlePath, "article.yml")

	err = r.fs.MkdirAll(r.fs.Join(articlePath, "images"), 0o777)
	runtime.Must(err)

	article := &model.Article{Name: name, Slug: slug}
	articles.Articles = append(articles.Articles, article)

	configFile, err := r.fs.Create(articlesConfig)
	runtime.Must(err)

	defer r.close(configFile)

	err = r.encoder.Encode(configFile, articles)
	runtime.Must(err)

	articleConfigFile, err := r.fs.Create(articleConfigPath)
	runtime.Must(err)

	defer r.close(articleConfigFile)

	err = r.encoder.Encode(articleConfigFile, article)
	runtime.Must(err)

	articleBodyPath := r.fs.Join(articlePath, "article.md")

	articleBodyFile, err := r.fs.Create(articleBodyPath)
	runtime.Must(err)

	defer r.close(articleBodyFile)

	return nil
}

// PublishArticle to the bucket.
func (r *S3Repository) PublishArticle(ctx context.Context, slug string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("repository: publish article", runtime.ConvertRecover(r))
		}
	}()

	ctx = meta.WithRequestID(ctx, meta.String(r.generator.Generate()))

	articlesPath, articlesConfig := r.configPath()
	r.uploadConfig(ctx, articlesConfig)

	articlePath := r.fs.Join(articlesPath, slug)
	r.uploadArticle(ctx, slug, articlePath)

	imagesPath := r.fs.Join(articlePath, "images")
	r.uploadImages(ctx, slug, imagesPath)

	return nil
}

// UnpublishArticle from the bucket.
func (r *S3Repository) UnpublishArticle(ctx context.Context, slug string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("repository: unpublish article", runtime.ConvertRecover(r))
		}
	}()

	ctx = meta.WithRequestID(ctx, meta.String(r.generator.Generate()))
	articles := r.articles(ctx)
	articlesPath, articlesConfig := r.configPath()
	articlePath := r.fs.Join(articlesPath, slug)

	r.delete(ctx, articlesPath, articlePath)

	err = r.fs.RemoveAll(articlePath)
	runtime.Must(err)

	r.deleteConfig(ctx, slug, articlesConfig, articles)

	return nil
}

func (r *S3Repository) uploadConfig(ctx context.Context, path string) {
	r.put(ctx, "articles.yml", mime.YAMLMediaType, path)
}

func (r *S3Repository) uploadArticle(ctx context.Context, slug, path string) {
	configPath := r.fs.Join(path, "article.yml")
	r.put(ctx, r.fs.Join(slug, "article.yml"), mime.YAMLMediaType, configPath)

	bodyPath := r.fs.Join(path, "article.md")
	r.put(ctx, r.fs.Join(slug, "article.md"), mime.MarkdownMediaType, bodyPath)
}

func (r *S3Repository) uploadImages(ctx context.Context, slug, path string) {
	_ = r.fs.WalkDir(path, func(path string, info os.DirEntry, err error) error {
		runtime.Must(err)

		if info.IsDir() {
			return nil
		}

		r.put(ctx, r.fs.Join(slug, "images", r.fs.Base(path)), mime.JPEGMediaType, path)

		return nil
	})
}

func (r *S3Repository) deleteConfig(ctx context.Context, slug, path string, articles *model.Articles) {
	articles.Articles = slices.DeleteFunc(articles.Articles, func(a *model.Article) bool { return a.Slug == slug })

	file, err := r.fs.Create(path)
	runtime.Must(err)

	defer r.close(file)

	err = r.encoder.Encode(file, articles)
	runtime.Must(err)

	r.uploadConfig(ctx, path)
}

func (r *S3Repository) delete(ctx context.Context, base, path string) {
	_ = r.fs.WalkDir(path, func(path string, info os.DirEntry, err error) error {
		runtime.Must(err)

		if info.IsDir() {
			return nil
		}

		rel, err := r.fs.Rel(base, path)
		runtime.Must(err)

		input := &s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    ptr.Value(rel),
		}

		_, err = r.s3.DeleteObject(ctx, input)
		runtime.Must(err)

		return nil
	})
}

func (r *S3Repository) put(ctx context.Context, path, contentType, body string) {
	file, err := r.fs.Open(body)
	runtime.Must(err)

	defer r.close(file)

	input := &s3.PutObjectInput{
		Bucket:      bucket,
		Key:         ptr.Value(path),
		Body:        file,
		ContentType: ptr.Value(contentType),
	}

	_, err = r.s3.PutObject(ctx, input)
	runtime.Must(err)
}

func (r *S3Repository) articles(ctx context.Context) *model.Articles {
	site := &model.Articles{}
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    ptr.Value("articles.yml"),
	}

	out, err := r.s3.GetObject(ctx, input)
	if err != nil {
		if s3.IsNotFound(err) {
			return site
		}

		runtime.Must(err)

		return nil
	}

	err = r.encoder.Decode(out.Body, site)
	runtime.Must(err)

	return site
}

func (r *S3Repository) close(closer io.Closer) {
	err := closer.Close()
	runtime.Must(err)
}

func (r *S3Repository) configPath() (string, string) {
	articlesPath := r.fs.Join(r.config.GetPath(r.fs), "articles")
	articlesConfig := r.fs.Join(articlesPath, "articles.yml")

	return articlesPath, articlesConfig
}
