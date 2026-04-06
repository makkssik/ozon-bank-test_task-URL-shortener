package redirect

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/presentation/http/handlers/redirect/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	type testCase struct {
		name             string
		alias            string
		mockBehavior     func(m *mocks.URLProvider)
		expectedStatus   int
		expectedLocation string
	}

	cases := []testCase{
		{
			name:  "Successful redirect",
			alias: "AbCdEfGhIj",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetOriginal", mock.Anything, "AbCdEfGhIj").
					Return(entities.RestoreURL("https://github.com", "AbCdEfGhIj"), nil).Once()
			},
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://github.com",
		},
		{
			name:  "URL not found",
			alias: "notfound12",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetOriginal", mock.Anything, "notfound12").
					Return((*entities.URL)(nil), nil).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "Internal server error",
			alias: "error12345",
			mockBehavior: func(m *mocks.URLProvider) {
				m.On("GetOriginal", mock.Anything, "error12345").
					Return((*entities.URL)(nil), context.DeadlineExceeded).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockProvider := mocks.NewURLProvider(t)
			tc.mockBehavior(mockProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewRedirectHandler(logger, mockProvider)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.alias, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusFound {
				require.Equal(t, tc.expectedLocation, rr.Header().Get("Location"))
			}
		})
	}
}
