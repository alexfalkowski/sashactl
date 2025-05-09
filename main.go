package main

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/sashactl/internal/cmd/new"
	"github.com/alexfalkowski/sashactl/internal/cmd/publish"
	"github.com/alexfalkowski/sashactl/internal/cmd/unpublish"
)

func main() {
	command().ExitOnError(context.Background())
}

func command() *cmd.Command {
	command := cmd.New(env.NewName(), env.NewVersion())

	new.Register(command)
	publish.Register(command)
	unpublish.Register(command)

	return command
}
