package entities

import (
	"testing"
	"url-shortener/internal/domain/exceptions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{
			name: "Success: lowercase scheme and host",
			raw:  "HTTPS://GOOGLE.COM/Path/",
			want: "https://google.com/Path",
		},
		{
			name: "Success: trim trailing slash",
			raw:  "https://example.com/test/",
			want: "https://example.com/test",
		},
		{
			name:    "Error: invalid URL",
			raw:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.raw)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		wantErr error
	}{
		{
			name:    "Valid alias",
			alias:   "aBcDeF123_",
			wantErr: nil,
		},
		{
			name:    "Invalid length (too short)",
			alias:   "short",
			wantErr: exceptions.ErrInvalidChars,
		},
		{
			name:    "Invalid length (too long)",
			alias:   "too_long_alias_123",
			wantErr: exceptions.ErrInvalidChars,
		},
		{
			name:    "Invalid characters (symbols)",
			alias:   "abcDEF12@#",
			wantErr: exceptions.ErrInvalidChars,
		},
		{
			name:    "Invalid characters (russian)",
			alias:   "привет_мир",
			wantErr: exceptions.ErrInvalidChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.alias)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewURL(t *testing.T) {
	original := "https://google.com"
	url := NewURL(original)

	require.NotNil(t, url)
	assert.Equal(t, original, url.Original)
	assert.Len(t, url.Short, 10)
	assert.NoError(t, Validate(url.Short))
}
