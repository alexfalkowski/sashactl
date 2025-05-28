package delete

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	se "github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/module"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/cmd/errors"
	"github.com/alexfalkowski/sashactl/internal/config"
	"go.uber.org/fx"
)

// Register for delete.
func Register(command cli.Commander) {
	cmd := command.AddClient("delete", "Delete a new article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, cli.Module, Module,
	)
	cmd.AddInput("")
	cmd.StringP("slug", "s", "", "slug of the article")
}

// Params for delete.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *flag.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Delete a new article.
func Delete(params Params) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slug, _ := params.FlagSet.GetString("slug")
			if strings.IsEmpty(slug) {
				return errors.ErrNoSlug
			}

			if err := params.Repository.DeleteArticle(ctx, slug); err != nil {
				return se.Prefix("delete: created article", err)
			}

			params.Logger.Info("deleted article", slog.String("slug", slug))

			return nil
		},
	})
}
