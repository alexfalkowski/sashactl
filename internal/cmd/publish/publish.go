package publish

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/config"
	"go.uber.org/fx"
)

// Register for publish.
func Register(command *cmd.Command) {
	flags := command.AddClient("publish", "Publish the article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, Module, cmd.Module,
	)
	flags.AddInput("")
	flags.StringP("slug", "s", "", "slug of the article")
}

// Params for publish.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *cmd.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Publish the created article.
func Publish(params Params) {
	cmd.Start(params.Lifecycle, func(ctx context.Context) error {
		slug, _ := params.FlagSet.GetString("slug")

		if err := params.Repository.PublishArticle(ctx, slug); err != nil {
			return errors.Prefix("publish: existing article", err)
		}

		params.Logger.Info("published article", slog.String("slug", slug))

		return nil
	})
}
