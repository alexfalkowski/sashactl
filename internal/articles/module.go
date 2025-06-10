package articles

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/sashactl/internal/articles/repository"
)

// Module for fx.
var Module = di.Module(
	repository.Module,
)
