package s3

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	ac "github.com/alexfalkowski/sashactl/internal/aws/config"
	"github.com/alexfalkowski/sashactl/internal/aws/endpoint"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type (
	// Client is an alias for s3.Client.
	Client = s3.Client

	// DeleteObjectInput is an alias for s3.DeleteObjectInput.
	DeleteObjectInput = s3.DeleteObjectInput

	// GetObjectInput is an alias for s3.GetObjectInput.
	GetObjectInput = s3.GetObjectInput

	// PutObjectInput is an alias for s3.PutObjectInput.
	PutObjectInput = s3.PutObjectInput
)

// IsNotFound for s3.
func IsNotFound(err error) bool {
	var noKeyErr *types.NoSuchKey

	return errors.As(err, &noKeyErr)
}

// ConfigParams for s3.
type ClientParams struct {
	di.In

	Tracer    *tracer.Tracer
	Meter     *metrics.Meter
	ID        id.Generator
	Endpoint  endpoint.Endpoint
	Config    *ac.Config
	Logger    *logger.Logger
	FS        *os.FS
	Limiter   *limiter.Limiter
	UserAgent env.UserAgent
}

// NewClient for s3.
func NewClient(params ClientParams) (*Client, error) {
	// As recommended by https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/.
	config.WithRequestChecksumCalculation(0)
	config.WithResponseChecksumValidation(0)

	accessKeyID, err := params.Config.GetAccessKeyID(params.FS)
	if err != nil {
		return nil, err
	}

	accessKeySecret, err := params.Config.GetAccessKeySecret(params.FS)
	if err != nil {
		return nil, err
	}

	httpClient, _ := http.NewClient(
		http.WithClientLogger(params.Logger), http.WithClientTracer(params.Tracer),
		http.WithClientMetrics(params.Meter), http.WithClientUserAgent(params.UserAgent),
		http.WithClientTimeout(params.Config.Timeout), http.WithClientID(params.ID),
		http.WithClientLimiter(params.Limiter),
	)

	opts := []func(*config.LoadOptions) error{
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				bytes.String(accessKeyID),
				bytes.String(accessKeySecret),
				"",
			),
		),
		config.WithRegion(params.Config.Region),
		config.WithHTTPClient(httpClient),
		config.WithRetryMaxAttempts(int(params.Config.Retry.Attempts)), //nolint:gosec
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true

		if params.Endpoint.IsSet() {
			o.BaseEndpoint = aws.String(params.Endpoint.String())
		}
	})

	return client, err
}
