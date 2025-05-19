package main

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/sashactl/internal/cmd/delete"
	"github.com/alexfalkowski/sashactl/internal/cmd/new"
	"github.com/alexfalkowski/sashactl/internal/cmd/publish"
	"github.com/alexfalkowski/sashactl/internal/cmd/unpublish"
)

var app = cmd.NewApplication(func(command *cmd.Command) {
	delete.Register(command)
	new.Register(command)
	publish.Register(command)
	unpublish.Register(command)
})

func main() {
	app.ExitOnError(context.Background())
}
