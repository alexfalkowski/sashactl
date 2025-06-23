package new

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

// Register for new.
func Register(command cli.Commander) {
	cmd := command.AddClient("new", "Create a new article", Module)

	cmd.AddInput("")
	cmd.StringP("name", "n", "", "name of the article")
}

// Params for new.
type Params struct {
	di.In

	Lifecycle  di.Lifecycle
	Logger     *logger.Logger
	FlagSet    *flag.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// New article to be created.
func New(params Params) {
	params.Lifecycle.Append(di.Hook{
		OnStart: func(ctx context.Context) error {
			name, _ := params.FlagSet.GetString("name")
			if strings.IsEmpty(name) {
				return errors.Prefix("new", errors.ErrNoName)
			}

			if err := params.Repository.NewArticle(ctx, name); err != nil {
				return errors.Prefix("new: create article", err)
			}

			params.Logger.Info("created article", logger.String("name", name))

			return nil
		},
	})
}
