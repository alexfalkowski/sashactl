package articles

import (
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"github.com/alexfalkowski/sashactl/internal/aws"
	"github.com/alexfalkowski/sashactl/internal/slug"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	http.Module,
	slug.Module,
	aws.Module,
	repository.Module,
)
