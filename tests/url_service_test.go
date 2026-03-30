package tests

import (
	"context"
	"testing"
	"url-shortener/internal/application/abstractions/repositories/mocks"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/application/services"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestURLService_ShortenURL_ShouldGenerateAndSave(t *testing.T) {
	mockRepo := mocks.NewURLRepository(t)
	svc := services.NewURLService(mockRepo)

	originalURL := "https://google.com"
	normalizedURL := "https://google.com"

	mockRepo.On("GetByOriginal", mock.Anything, normalizedURL).Return(nil, exceptions.ErrURLNotFound).Once()

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.URL")).Return(nil).Once()

	res, err := svc.ShortenURL(context.Background(), operations.ShortenRequest{
		OriginalURL: originalURL,
	})

	require.NoError(t, err)
	require.Len(t, res.ShortURL, 10)
}

func TestURLService_GetOriginal_ShouldReturnURL(t *testing.T) {
	mockRepo := mocks.NewURLRepository(t)
	svc := services.NewURLService(mockRepo)

	shortURL := "1234567890"
	expectedURL := &entities.URL{
		Original: "https://vk.com",
		Short:    shortURL,
	}

	mockRepo.On("GetByShort", mock.Anything, shortURL).Return(expectedURL, nil).Once()

	res, err := svc.GetOriginal(context.Background(), operations.GetOriginalRequest{
		ShortURL: shortURL,
	})

	require.NoError(t, err)
	require.Equal(t, "https://vk.com", res.OriginalURL)
}
