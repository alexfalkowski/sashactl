package unpublish

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/cmd/errors"
	"github.com/alexfalkowski/sashactl/internal/config"
)

// Register for unpublish.
func Register(command cli.Commander) {
	cmd := command.AddClient("unpublish", "Unpublish an article", Module)

	cmd.AddInput("")
	cmd.StringP("slug", "s", "", "slug of the article")
}

// Params for unpublish.
type Params struct {
	di.In

	Lifecycle  di.Lifecycle
	Logger     *logger.Logger
	FlagSet    *flag.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Unpublish a article.
func Unpublish(params Params) {
	params.Lifecycle.Append(di.Hook{
		OnStart: func(ctx context.Context) error {
			slug, _ := params.FlagSet.GetString("slug")
			if strings.IsEmpty(slug) {
				return errors.Prefix("unpublish", errors.ErrNoSlug)
			}

			if err := params.Repository.UnpublishArticle(ctx, slug); err != nil {
				return errors.Prefix("unpublish: existing article", err)
			}

			params.Logger.Info("unpublished article", slog.String("slug", slug))

			return nil
		},
	})
}
