package repositories

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveOrGet(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	query := `
		WITH inserted AS (
			INSERT INTO urls (original, short) 
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			RETURNING short
		)
		SELECT short FROM inserted
		UNION ALL
		SELECT short FROM urls WHERE original = $1
		LIMIT 1;
	`

	var existingShort string

	err := s.db.QueryRowContext(ctx, query, url.Original, url.Short).Scan(&existingShort)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrAliasCollision
		}
		return nil, err
	}

	return entities.RestoreURL(url.Original, existingShort), nil
}

func (s *Storage) GetByShort(ctx context.Context, short string) (*entities.URL, error) {
	query := `SELECT original FROM urls WHERE short = $1`

	var originalURL string

	err := s.db.QueryRowContext(ctx, query, short).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return entities.RestoreURL(originalURL, short), nil
}
