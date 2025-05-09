package unpublish

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/encoding/yaml"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/cmd/errors"
	"github.com/alexfalkowski/sashactl/internal/config"
	"go.uber.org/fx"
)

// Register for unpublish.
func Register(command *cmd.Command) {
	flags := command.AddClient("unpublish", "Unpublish an article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, Module, cmd.Module,
	)
	flags.AddInput("")
	flags.StringP("slug", "s", "", "slug of the article")
}

// Params for unpublish.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *cmd.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Unpublish a article.
func Unpublish(params Params) {
	cmd.Start(params.Lifecycle, func(ctx context.Context) error {
		slug, _ := params.FlagSet.GetString("slug")
		if strings.IsEmpty(slug) {
			return errors.ErrNoSlug
		}

		if err := params.Repository.UnpublishArticle(ctx, slug); err != nil {
			return se.Prefix("unpublish: existing article", err)
		}

		params.Logger.Info("unpublished article", slog.String("slug", slug))

		return nil
	})
}
