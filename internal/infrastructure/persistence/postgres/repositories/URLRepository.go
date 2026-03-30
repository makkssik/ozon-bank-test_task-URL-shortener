package repositories

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"

	"github.com/lib/pq"
)

const pgErrorCodeUniqueViolation = "23505"

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Save(ctx context.Context, url *entities.URL) error {
	query := `INSERT INTO urls (original, short) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, url.Original, url.Short)
	if err != nil {
		if pqErr, ok := errors.AsType[*pq.Error](err); ok && pqErr.Code == pgErrorCodeUniqueViolation {
			return exceptions.ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (s *Storage) GetByShort(ctx context.Context, short string) (*entities.URL, error) {
	query := `SELECT original, short FROM urls WHERE short = $1`

	var originalURL, shortURL string
	err := s.db.QueryRowContext(ctx, query, short).Scan(&originalURL, &shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrURLNotFound
		}
		return nil, err
	}

	return &entities.URL{Original: originalURL, Short: shortURL}, nil
}

func (s *Storage) GetByOriginal(ctx context.Context, original string) (*entities.URL, error) {
	query := `SELECT original, short FROM urls WHERE original = $1`

	var originalURL, shortURL string
	err := s.db.QueryRowContext(ctx, query, original).Scan(&originalURL, &shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrURLNotFound
		}
		return nil, err
	}
	return &entities.URL{Original: originalURL, Short: shortURL}, nil
}
