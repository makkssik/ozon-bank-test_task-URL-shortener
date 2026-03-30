package repositories

import (
	"context"
	"url-shortener/internal/domain/entities"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLRepository --output=./mocks --outpkg=mocks
type URLRepository interface {
	Save(ctx context.Context, url *entities.URL) error

	GetByShort(ctx context.Context, short string) (*entities.URL, error)

	GetByOriginal(ctx context.Context, original string) (*entities.URL, error)
}
