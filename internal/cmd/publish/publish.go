package publish

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/cmd/errors"
	"github.com/alexfalkowski/sashactl/internal/config"
)

// Register for publish.
func Register(command cli.Commander) {
	cmd := command.AddClient("publish", "Publish the article", Module)

	cmd.AddInput("")
	cmd.StringP("slug", "s", "", "slug of the article")
}

// Params for publish.
type Params struct {
	di.In

	Lifecycle  di.Lifecycle
	Logger     *logger.Logger
	FlagSet    *flag.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Publish the created article.
func Publish(params Params) {
	params.Lifecycle.Append(di.Hook{
		OnStart: func(ctx context.Context) error {
			slug, _ := params.FlagSet.GetString("slug")
			if strings.IsEmpty(slug) {
				return errors.Prefix("publish", errors.ErrNoSlug)
			}

			if err := params.Repository.PublishArticle(ctx, slug); err != nil {
				return errors.Prefix("publish: created article", err)
			}

			params.Logger.Info("published article", slog.String("slug", slug))

			return nil
		},
	})
}
