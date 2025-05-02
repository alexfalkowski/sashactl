package new

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/model"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/config"
	"github.com/gosimple/slug"
	"go.uber.org/fx"
)

// Params for config.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *cmd.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Start the creating an article.
func Start(params Params) {
	cmd.Start(params.Lifecycle, func(ctx context.Context) error {
		articles, err := params.Repository.GetArticles(ctx)
		if err != nil {
			return errors.Prefix("new: get articles", err)
		}

		name, _ := params.FlagSet.GetString("name")
		slug := slug.Make(name)

		articlesDir := filepath.Join(params.Config.Articles.Path, "articles")
		articlesConfig := filepath.Join(articlesDir, "articles.yml")
		articleDir := filepath.Join(articlesDir, slug)
		articleConfig := filepath.Join(articleDir, "article.yml")

		if err := os.MkdirAll(filepath.Join(articleDir, "images"), 0o777); err != nil {
			return errors.Prefix("new: mkdir", err)
		}

		article := &model.Article{Name: name, Slug: slug}
		articles.Articles = append(articles.Articles, article)

		configFile, err := os.Create(articlesConfig)
		if err != nil {
			return errors.Prefix("new: create articles", err)
		}

		if err := params.Encoder.Encode(configFile, articles); err != nil {
			return errors.Prefix("new: encode articles", err)
		}

		article.Body = "Add my story!"
		article.Images = []*model.Image{
			{Name: "dummy", Description: "Add me!"},
		}

		articleFile, err := os.Create(articleConfig)
		if err != nil {
			return errors.Prefix("new: create article", err)
		}

		if err := params.Encoder.Encode(articleFile, article); err != nil {
			return errors.Prefix("new: encode article", err)
		}

		params.Logger.Info("created article", slog.String("name", name))

		return nil
	})
}
