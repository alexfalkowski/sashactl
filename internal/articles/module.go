package articles

import (
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	repository.Module,
)
