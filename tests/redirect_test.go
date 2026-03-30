package tests

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/application/contracts/url/mocks"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/presentation/http/handlers/redirect"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRedirectTest(t *testing.T) (*mocks.Service, *httptest.ResponseRecorder, http.HandlerFunc) {
	mockService := mocks.NewService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := redirect.NewRedirectHandler(logger, mockService)
	rr := httptest.NewRecorder()

	return mockService, rr, handler
}

func createRequestWithAlias(t *testing.T, alias string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("alias", alias)
	req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+alias, nil)
	require.NoError(t, err)

	return req
}

func TestRedirectHandler_ShouldRedirect_WhenAliasExists(t *testing.T) {
	mockService, rr, handler := setupRedirectTest(t)
	shortURL := "1234567890"
	originalURL := "https://github.com"

	mockService.On("GetOriginal", mock.Anything, operations.GetOriginalRequest{
		ShortURL: shortURL,
	}).Return(operations.GetOriginalResponse{OriginalURL: originalURL}, nil).Once()

	req := createRequestWithAlias(t, shortURL)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusFound, rr.Code)
	require.Equal(t, originalURL, rr.Header().Get("Location"))
}

func TestRedirectHandler_ShouldReturnNotFound_WhenAliasDoesNotExist(t *testing.T) {
	mockService, rr, handler := setupRedirectTest(t)
	shortURL := "notfound12"

	mockService.On("GetOriginal", mock.Anything, operations.GetOriginalRequest{
		ShortURL: shortURL,
	}).Return(operations.GetOriginalResponse{}, exceptions.ErrURLNotFound).Once()

	req := createRequestWithAlias(t, shortURL)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}
