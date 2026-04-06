package save

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/presentation/http/handlers/save/mocks"
	"url-shortener/internal/presentation/http/response"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockBehavior   func(m *mocks.URLShortener)
		expectedStatus int
		expectedAlias  string
		expectedError  string
	}

	cases := []testCase{
		{
			name: "Success: URL shortened",
			url:  "https://google.com",
			mockBehavior: func(m *mocks.URLShortener) {
				m.On("ShortenURL", mock.Anything, "https://google.com").
					Return(entities.RestoreURL("https://google.com", "AbCdEfGhIj"), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedAlias:  "AbCdEfGhIj",
		},
		{
			name:           "Error: Invalid URL format",
			url:            "not-a-valid-url",
			mockBehavior:   func(m *mocks.URLShortener) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field URL is invalid url",
		},
		{
			name:           "Error: Empty URL",
			url:            "",
			mockBehavior:   func(m *mocks.URLShortener) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field URL is required field",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockShortener := mocks.NewURLShortener(t)
			tc.mockBehavior(mockShortener)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewSaveHandler(logger, mockShortener)

			reqBody, _ := json.Marshal(Request{URL: tc.url})
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			var resp Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.expectedStatus == http.StatusCreated {
				require.Equal(t, response.StatusOk, resp.Status)
				require.Equal(t, tc.expectedAlias, resp.Alias)
			} else {
				require.Equal(t, response.StatusError, resp.Status)
				require.Contains(t, resp.Error, tc.expectedError)
			}
		})
	}
}
