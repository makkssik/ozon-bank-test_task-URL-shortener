package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"url-shortener/internal/application/abstractions/repositories"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
)

type URLService struct {
	repository repositories.URLRepository
}

func NewURLService(repository repositories.URLRepository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)

	if len(parsedURL.Path) > 1 {
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	}

	return parsedURL.String(), nil
}

func (s *URLService) ShortenURL(ctx context.Context, request operations.ShortenRequest) (operations.ShortenResponse, error) {
	normalizedURL, err := normalizeURL(request.OriginalURL)
	if err != nil {
		return operations.ShortenResponse{}, fmt.Errorf("invalid url format: %w", err)
	}

	existingURL, err := s.repository.GetByOriginal(ctx, normalizedURL)
	if err == nil {

		return operations.ShortenResponse{ShortURL: existingURL.Short}, nil
	}

	if !errors.Is(err, exceptions.ErrURLNotFound) {
		return operations.ShortenResponse{}, fmt.Errorf("failed to check existing url: %w", err)
	}

	const maxRetries = 3
	var newURL *entities.URL

	for i := 0; i < maxRetries; i++ {
		shortPath := getRandomString(10)
		newURL, err = entities.NewURL(normalizedURL, shortPath)

		if err != nil {
			return operations.ShortenResponse{}, fmt.Errorf("invalid url generated: %w", err)
		}

		err = s.repository.Save(ctx, newURL)
		if err == nil {
			return operations.ShortenResponse{ShortURL: newURL.Short}, nil
		}

		if errors.Is(err, exceptions.ErrAlreadyExists) {
			existing, checkErr := s.repository.GetByOriginal(ctx, normalizedURL)
			if checkErr == nil {
				return operations.ShortenResponse{ShortURL: existing.Short}, nil
			}
			continue
		}
		return operations.ShortenResponse{}, fmt.Errorf("failed to save url: %w", err)
	}

	return operations.ShortenResponse{}, errors.New("failed to generate url after retries")
}

func (s *URLService) GetOriginal(ctx context.Context, request operations.GetOriginalRequest) (operations.GetOriginalResponse, error) {
	domainURL, err := s.repository.GetByShort(ctx, request.ShortURL)
	if err != nil {
		if errors.Is(err, exceptions.ErrURLNotFound) {
			return operations.GetOriginalResponse{}, err
		}
		return operations.GetOriginalResponse{}, fmt.Errorf("failed to get url: %w", err)
	}

	return operations.GetOriginalResponse{OriginalURL: domainURL.Original}, nil
}

func getRandomString(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	bytes := make([]byte, n)

	for i := range bytes {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		bytes[i] = alphabet[num.Int64()]
	}

	return string(bytes)
}
