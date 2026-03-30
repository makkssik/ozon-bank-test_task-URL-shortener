package url

import (
	"context"
	"url-shortener/internal/application/contracts/url/operations"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=Service --output=./mocks --outpkg=mocks
type Service interface {
	ShortenURL(ctx context.Context, request operations.ShortenRequest) (operations.ShortenResponse, error)

	GetOriginal(ctx context.Context, request operations.GetOriginalRequest) (operations.GetOriginalResponse, error)
}
