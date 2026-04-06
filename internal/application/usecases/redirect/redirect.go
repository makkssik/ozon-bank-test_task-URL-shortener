package redirect

import (
	"context"
	"fmt"
	"url-shortener/internal/domain/entities"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLProvider --output=./mocks --outpkg=mocks --filename=url_provider.go
type URLProvider interface {
	GetByShort(ctx context.Context, short string) (*entities.URL, error)
}

type Service struct {
	provider URLProvider
}

func NewService(provider URLProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) GetOriginal(ctx context.Context, shortURL string) (*entities.URL, error) {
	if err := entities.Validate(shortURL); err != nil {
		return nil, fmt.Errorf("invalid alias format: %w", err)
	}

	urlEntity, err := s.provider.GetByShort(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get url from db: %w", err)
	}

	if urlEntity == nil {
		return nil, nil
	}

	return urlEntity, nil
}
