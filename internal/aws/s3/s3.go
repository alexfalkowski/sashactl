package s3

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http"
	ac "github.com/alexfalkowski/sashactl/internal/aws/config"
	"github.com/alexfalkowski/sashactl/internal/aws/endpoint"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/fx"
)

// ConfigParams for s3.
type ClientParams struct {
	fx.In
	Tracer    *tracer.Tracer
	Meter     *metrics.Meter
	ID        id.Generator
	Endpoint  endpoint.Endpoint
	Config    *ac.Config
	Logger    *logger.Logger
	UserAgent env.UserAgent
}

// NewClient for s3.
func NewClient(params ClientParams) (*s3.Client, error) {
	// As recommended by https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/.
	config.WithRequestChecksumCalculation(0)
	config.WithResponseChecksumValidation(0)

	accessKeyID, err := params.Config.GetAccessKeyID()
	if err != nil {
		return nil, err
	}

	accessKeySecret, err := params.Config.GetAccessKeySecret()
	if err != nil {
		return nil, err
	}

	httpClient, _ := http.NewClient(
		http.WithClientLogger(params.Logger), http.WithClientTracer(params.Tracer),
		http.WithClientMetrics(params.Meter), http.WithClientUserAgent(params.UserAgent),
		http.WithClientTimeout(params.Config.Timeout), http.WithClientID(params.ID),
	)

	opts := []func(*config.LoadOptions) error{
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(string(accessKeyID), string(accessKeySecret), "")),
		config.WithRegion(params.Config.Region),
		config.WithHTTPClient(httpClient),
		config.WithRetryMaxAttempts(int(params.Config.Retry.Attempts)), //nolint:gosec
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true

		if params.Endpoint.IsSet() {
			o.BaseEndpoint = aws.String(params.Endpoint.String())
		}
	})

	return s3Client, err
}
