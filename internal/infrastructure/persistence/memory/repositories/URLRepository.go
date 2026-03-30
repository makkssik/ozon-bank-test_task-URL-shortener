package repositories

import (
	"context"
	"sync"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
)

type Storage struct {
	mu sync.RWMutex

	shortToUrl    map[string]*entities.URL
	originalToUrl map[string]*entities.URL
}

func NewStorage() *Storage {
	return &Storage{
		shortToUrl:    make(map[string]*entities.URL),
		originalToUrl: make(map[string]*entities.URL),
	}
}

func (s *Storage) Save(ctx context.Context, url *entities.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.shortToUrl[url.Short]; exists {
		return exceptions.ErrAlreadyExists
	}
	if _, exists := s.originalToUrl[url.Original]; exists {
		return exceptions.ErrAlreadyExists
	}

	s.shortToUrl[url.Short] = url
	s.originalToUrl[url.Original] = url

	return nil
}

func (s *Storage) GetByShort(ctx context.Context, short string) (*entities.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.shortToUrl[short]
	if !exists {
		return nil, exceptions.ErrURLNotFound
	}

	return url, nil
}

func (s *Storage) GetByOriginal(ctx context.Context, original string) (*entities.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.originalToUrl[original]
	if !exists {
		return nil, exceptions.ErrURLNotFound
	}

	return url, nil
}
