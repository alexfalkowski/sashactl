package delete

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

// Register for delete.
func Register(command *cmd.Command) {
	flags := command.AddClient("delete", "Delete a new article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, Module, cmd.Module,
	)
	flags.AddInput("")
	flags.StringP("slug", "s", "", "slug of the article")
}

// Params for delete.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *cmd.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Delete a new article.
func Delete(params Params) {
	cmd.Start(params.Lifecycle, func(ctx context.Context) error {
		slug, _ := params.FlagSet.GetString("slug")
		if strings.IsEmpty(slug) {
			return errors.ErrNoSlug
		}

		if err := params.Repository.DeleteArticle(ctx, slug); err != nil {
			return se.Prefix("delete: created article", err)
		}

		params.Logger.Info("deleted article", slog.String("slug", slug))

		return nil
	})
}
