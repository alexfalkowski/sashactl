package unpublish

import (
	"github.com/alexfalkowski/sashactl/internal/articles"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	articles.Module,
	fx.Invoke(Unpublish),
)
