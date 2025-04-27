package cmd

import (
	"github.com/alexfalkowski/go-client-template/internal/cmd/client"
	"github.com/alexfalkowski/go-client-template/internal/config"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/telemetry"
)

// RegisterClient for cmd.
func RegisterClient(command *cmd.Command) {
	flags := command.AddClient("client", "Start client",
		module.Module, feature.Module, telemetry.Module,
		config.Module, client.Module, cmd.Module,
	)
	flags.AddInput("")
}
