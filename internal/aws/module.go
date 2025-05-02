package aws

import (
	aws "github.com/alexfalkowski/sashactl/internal/aws/endpoint"
	"github.com/alexfalkowski/sashactl/internal/aws/s3"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(aws.NewEndpoint),
	fx.Provide(s3.NewClient),
)
