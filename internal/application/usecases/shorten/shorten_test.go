package shorten

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/application/usecases/shorten/mocks"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_ShortenURL(t *testing.T) {
	type testCase struct {
		name          string
		originalURL   string
		mockBehavior  func(m *mocks.URLSaver)
		expectedError string
	}

	cases := []testCase{
		{
			name:        "URL successfully saved",
			originalURL: "https://google.com",
			mockBehavior: func(m *mocks.URLSaver) {
				m.On("SaveOrGet", mock.Anything, mock.AnythingOfType("*entities.URL")).
					Return(entities.RestoreURL("https://google.com", "AbCdEfGhIj"), nil).Once()
			},
			expectedError: "",
		},
		{
			name:          "Invalid URL format",
			originalURL:   "not-a-url",
			mockBehavior:  func(m *mocks.URLSaver) {},
			expectedError: "invalid url format",
		},
		{
			name:        "Max retries exceeded",
			originalURL: "https://google.com",
			mockBehavior: func(m *mocks.URLSaver) {
				m.On("SaveOrGet", mock.Anything, mock.AnythingOfType("*entities.URL")).
					Return((*entities.URL)(nil), exceptions.ErrAliasCollision).Times(3)
			},
			expectedError: "failed to generate unique short url after max retries",
		},
		{
			name:        "Database fatal error",
			originalURL: "https://google.com",
			mockBehavior: func(m *mocks.URLSaver) {
				m.On("SaveOrGet", mock.Anything, mock.AnythingOfType("*entities.URL")).
					Return((*entities.URL)(nil), errors.New("db connection lost")).Once()
			},
			expectedError: "failed to save url: db connection lost",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			saverMock := mocks.NewURLSaver(t)
			tc.mockBehavior(saverMock)

			service := NewService(saverMock)

			urlEntity, err := service.ShortenURL(context.Background(), tc.originalURL)

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
				require.Nil(t, urlEntity)
			} else {
				require.NoError(t, err)
				require.NotNil(t, urlEntity)
				require.NotEmpty(t, urlEntity.Short)
			}
		})
	}
}
