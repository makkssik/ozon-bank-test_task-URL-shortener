package repositories

import (
	"context"
	"sync"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
)

type Storage struct {
	mu sync.RWMutex

	shortToOriginal map[string]string
	originalToShort map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

func (s *Storage) SaveOrGet(ctx context.Context, e *entities.URL) (*entities.URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existingShort, exists := s.originalToShort[e.Original]; exists {
		return entities.RestoreURL(e.Original, existingShort), nil
	}

	if _, exists := s.shortToOriginal[e.Short]; exists {
		return nil, exceptions.ErrAliasCollision
	}

	s.shortToOriginal[e.Short] = e.Original
	s.originalToShort[e.Original] = e.Short

	return e, nil
}

func (s *Storage) GetByShort(ctx context.Context, short string) (*entities.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	original, exists := s.shortToOriginal[short]
	if !exists {
		return nil, nil
	}

	return entities.RestoreURL(original, short), nil
}
