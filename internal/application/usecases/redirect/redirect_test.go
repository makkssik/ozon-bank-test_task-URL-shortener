package redirect

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/application/usecases/redirect/mocks"
	"url-shortener/internal/domain/entities"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetOriginal(t *testing.T) {
	type testCase struct {
		name          string
		shortURL      string
		mockBehavior  func(m *mocks.URLProvider)
		expectedURL   string
		expectedError string
	}

	cases := []testCase{
		{
			name:     "URL found",
			shortURL: "AbCdEfGhIj",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetByShort", mock.Anything, "AbCdEfGhIj").
					Return(entities.RestoreURL("https://google.com", "AbCdEfGhIj"), nil).Once()
			},
			expectedURL:   "https://google.com",
			expectedError: "",
		},
		{
			name:          "Invalid alias format (Domain validation)",
			shortURL:      "invalid_!@",
			mockBehavior:  func(m *mocks.URLProvider) {},
			expectedURL:   "",
			expectedError: "invalid alias format",
		},
		{
			name:     "URL not found in DB",
			shortURL: "AbCdEfGhIj",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetByShort", mock.Anything, "AbCdEfGhIj").
					Return((*entities.URL)(nil), nil).Once()
			},
			expectedURL:   "",
			expectedError: "",
		},
		{
			name:     "Database error",
			shortURL: "AbCdEfGhIj",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetByShort", mock.Anything, "AbCdEfGhIj").
					Return((*entities.URL)(nil), errors.New("timeout")).Once()
			},
			expectedURL:   "",
			expectedError: "failed to get url from db: timeout",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			providerMock := mocks.NewURLProvider(t)
			tc.mockBehavior(providerMock)

			service := NewService(providerMock)

			urlEntity, err := service.GetOriginal(context.Background(), tc.shortURL)

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
				require.Nil(t, urlEntity)
			} else {
				require.NoError(t, err)
				if tc.expectedURL == "" {
					require.Nil(t, urlEntity)
				} else {
					require.NotNil(t, urlEntity)
					require.Equal(t, tc.expectedURL, urlEntity.Original)
				}
			}
		})
	}
}
