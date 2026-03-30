package tests

import (
	"testing"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"

	"github.com/stretchr/testify/require"
)

func TestNewURL_ShouldCreate_WhenValid(t *testing.T) {
	url, err := entities.NewURL("https://example.com", "aBcDeF123_")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url.Original)
	require.Equal(t, "aBcDeF123_", url.Short)
}

func TestNewURL_ShouldReturnError_WhenInvalidLength(t *testing.T) {
	_, err := entities.NewURL("https://example.com", "short")
	require.ErrorIs(t, err, exceptions.ErrInvalidLength)

	_, err = entities.NewURL("https://example.com", "too_long_alias_123")
	require.ErrorIs(t, err, exceptions.ErrInvalidLength)
}

func TestNewURL_ShouldReturnError_WhenInvalidChars(t *testing.T) {
	_, err := entities.NewURL("https://example.com", "aBcDeF12-!")
	require.ErrorIs(t, err, exceptions.ErrInvalidChars)

	_, err = entities.NewURL("https://example.com", "abcDEF12@#")
	require.ErrorIs(t, err, exceptions.ErrInvalidChars)
}
