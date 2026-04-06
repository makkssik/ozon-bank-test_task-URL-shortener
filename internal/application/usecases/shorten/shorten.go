package shorten

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLSaver --output=./mocks --outpkg=mocks --filename=url_saver.go
type URLSaver interface {
	SaveOrGet(ctx context.Context, e *entities.URL) (*entities.URL, error)
}

type Service struct {
	saver URLSaver
}

func NewService(saver URLSaver) *Service {
	return &Service{saver: saver}
}

func (s *Service) ShortenURL(ctx context.Context, originalURL string) (*entities.URL, error) {
	normalizedURL, err := entities.NormalizeURL(originalURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url format: %w", err)
	}

	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		urlEntity := entities.NewURL(normalizedURL)

		savedURL, err := s.saver.SaveOrGet(ctx, urlEntity)
		if err != nil {
			if errors.Is(err, exceptions.ErrAliasCollision) {
				continue
			}
			return nil, fmt.Errorf("failed to save url: %w", err)
		}

		return savedURL, nil
	}

	return nil, errors.New("failed to generate unique short url after max retries")
}
