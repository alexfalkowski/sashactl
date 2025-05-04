package main

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/sashactl/internal/cmd/delete"
	"github.com/alexfalkowski/sashactl/internal/cmd/new"
	"github.com/alexfalkowski/sashactl/internal/cmd/publish"
)

func main() {
	command().ExitOnError(context.Background())
}

func command() *cmd.Command {
	command := cmd.New(env.NewName(), env.NewVersion())

	new.Register(command)
	publish.Register(command)
	delete.Register(command)

	return command
}
