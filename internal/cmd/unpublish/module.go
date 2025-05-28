package unpublish

import (
	"github.com/alexfalkowski/go-service/v2/module"
	"github.com/alexfalkowski/sashactl/internal/articles"
	"github.com/alexfalkowski/sashactl/internal/aws"
	"github.com/alexfalkowski/sashactl/internal/config"
	"github.com/alexfalkowski/sashactl/internal/slug"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	module.Client,
	config.Module,
	slug.Module,
	aws.Module,
	articles.Module,
	fx.Invoke(Unpublish),
)
