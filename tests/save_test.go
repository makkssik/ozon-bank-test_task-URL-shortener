package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/application/contracts/url/mocks"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/presentation/http/handlers/save"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSaveTest(t *testing.T) (*mocks.Service, *httptest.ResponseRecorder, http.HandlerFunc) {
	mockService := mocks.NewService(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := save.NewSaveHandler(logger, mockService)
	rr := httptest.NewRecorder()

	return mockService, rr, handler
}

func TestSaveHandler_ShouldReturnCreated_WhenRequestIsValid(t *testing.T) {
	mockService, rr, handler := setupSaveTest(t)

	mockService.On("ShortenURL", mock.Anything, operations.ShortenRequest{
		OriginalURL: "https://google.com",
	}).Return(operations.ShortenResponse{ShortURL: "AbCdEfGhIj"}, nil).Once()

	reqBody := `{"url": "https://google.com"}`
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var respBody save.Response
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	require.NoError(t, err)
	require.Equal(t, "AbCdEfGhIj", respBody.Alias)
	require.Equal(t, "OK", respBody.Status)
}

func TestSaveHandler_ShouldReturnBadRequest_WhenURLIsInvalid(t *testing.T) {
	_, rr, handler := setupSaveTest(t)

	reqBody := `{"url": "not-a-valid-url"}`
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveHandler_ShouldReturnConflict_WhenURLAlreadyExists(t *testing.T) {
	mockService, rr, handler := setupSaveTest(t)

	mockService.On("ShortenURL", mock.Anything, operations.ShortenRequest{
		OriginalURL: "https://google.com",
	}).Return(operations.ShortenResponse{}, exceptions.ErrAlreadyExists).Once()

	reqBody := `{"url": "https://google.com"}`
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)

	var respBody save.Response
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	require.NoError(t, err)
	require.Equal(t, "Error", respBody.Status)
	require.Equal(t, "url already exists", respBody.Error)
}
