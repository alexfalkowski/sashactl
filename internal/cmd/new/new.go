package new

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

// Register for new.
func Register(command *cmd.Command) {
	flags := command.AddClient("new", "Create a new article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, Module, cmd.Module,
	)
	flags.AddInput("")
	flags.StringP("name", "n", "", "name of the article")
}

// Params for new.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *logger.Logger
	FlagSet    *cmd.FlagSet
	Config     *config.Config
	Encoder    *yaml.Encoder
	Repository repository.Repository
}

// New article to be created.
func New(params Params) {
	cmd.Start(params.Lifecycle, func(ctx context.Context) error {
		name, _ := params.FlagSet.GetString("name")
		if strings.IsEmpty(name) {
			return errors.ErrNoName
		}

		if err := params.Repository.NewArticle(ctx, name); err != nil {
			return se.Prefix("new: create article", err)
		}

		params.Logger.Info("created article", slog.String("name", name))

		return nil
	})
}
