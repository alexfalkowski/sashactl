package unpublish

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/module"
	"github.com/alexfalkowski/sashactl/internal/articles"
	"github.com/alexfalkowski/sashactl/internal/aws"
	"github.com/alexfalkowski/sashactl/internal/config"
	"github.com/alexfalkowski/sashactl/internal/slug"
)

// Module for fx.
var Module = di.Module(
	module.Client,
	slug.Module,
	aws.Module,
	config.Module,
	articles.Module,
	di.Register(Unpublish),
)
