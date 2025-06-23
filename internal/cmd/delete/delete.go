package delete

import (
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

// Register for delete.
func Register(command cli.Commander) {
	cmd := command.AddClient("delete", "Delete a new article", Module)

	cmd.AddInput("")
	cmd.StringP("slug", "s", "", "slug of the article")
}

// Params for delete.
type Params struct {
	di.In

	Lifecycle  di.Lifecycle
	Logger     *logger.Logger
	FlagSet    *flag.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// Delete a new article.
func Delete(params Params) {
	params.Lifecycle.Append(di.Hook{
		OnStart: func(ctx context.Context) error {
			slug, _ := params.FlagSet.GetString("slug")
			if strings.IsEmpty(slug) {
				return errors.Prefix("delete", errors.ErrNoSlug)
			}

			if err := params.Repository.DeleteArticle(ctx, slug); err != nil {
				return errors.Prefix("delete: created article", err)
			}

			params.Logger.Info("deleted article", logger.String("slug", slug))

			return nil
		},
	})
}
