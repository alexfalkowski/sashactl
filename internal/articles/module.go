package articles

import (
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/sashactl/internal/articles/client"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	http.Module,
	client.Module,
	repository.Module,
)
