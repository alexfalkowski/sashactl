package endpoint

import (
	"github.com/alexfalkowski/go-service/v2/strings"
	aws "github.com/alexfalkowski/sashactl/internal/aws/config"
)

// NewEndpoint for aws.
func NewEndpoint(cfg *aws.Config) Endpoint {
	return Endpoint(cfg.Address)
}

// Endpoint for aws.
type Endpoint string

// IsSet for aws.
func (e Endpoint) IsSet() bool {
	return !strings.IsEmpty(e.String())
}

// String conforms to fmt.Stringer.
func (e Endpoint) String() string {
	return string(e)
}
