package main

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/sashactl/internal/cmd/delete"
	"github.com/alexfalkowski/sashactl/internal/cmd/new"
	"github.com/alexfalkowski/sashactl/internal/cmd/publish"
	"github.com/alexfalkowski/sashactl/internal/cmd/unpublish"
)

func main() {
	command().ExitOnError(context.Background())
}

func command() *cmd.Command {
	fs := os.NewFS()
	command := cmd.New(env.NewName(fs), env.NewVersion())

	delete.Register(command)
	new.Register(command)
	publish.Register(command)
	unpublish.Register(command)

	return command
}
