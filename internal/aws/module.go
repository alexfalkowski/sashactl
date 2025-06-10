package aws

import (
	"github.com/alexfalkowski/go-service/v2/di"
	aws "github.com/alexfalkowski/sashactl/internal/aws/endpoint"
	"github.com/alexfalkowski/sashactl/internal/aws/s3"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(aws.NewEndpoint),
	di.Constructor(s3.NewClient),
)
