package cmd

import (
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/sashactl/internal/cmd/new"
	"github.com/alexfalkowski/sashactl/internal/config"
)

// RegisterClient for cmd.
func RegisterClient(command *cmd.Command) {
	flags := command.AddClient("new", "New article",
		module.Module, feature.Module, telemetry.Module,
		config.Module, new.Module, cmd.Module,
	)
	flags.AddInput("")
	flags.StringP("name", "n", "", "name of the article")
}
