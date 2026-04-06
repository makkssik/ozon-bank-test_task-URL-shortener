package repositories

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"url-shortener/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestStorage_BasicOperations(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	type testCase struct {
		name        string
		original    string
		short       string
		action      func(s *Storage, original, short string) (string, error)
		expectedVal string
		expectedErr error
	}

	tests := []testCase{
		{
			name:     "Save new URL",
			original: "https://google.com",
			short:    "googl12345",
			action: func(s *Storage, original, short string) (string, error) {
				ent, err := s.SaveOrGet(ctx, entities.RestoreURL(original, short))
				if ent != nil {
					return ent.Short, err
				}
				return "", err
			},
			expectedVal: "googl12345",
			expectedErr: nil,
		},
		{
			name:     "Get existing URL",
			original: "https://google.com",
			short:    "googl12345",
			action: func(s *Storage, original, short string) (string, error) {
				ent, err := s.GetByShort(ctx, short)
				if ent != nil {
					return ent.Original, err
				}
				return "", err
			},
			expectedVal: "https://google.com",
			expectedErr: nil,
		},
		{
			name:     "Save duplicate original",
			original: "https://google.com",
			short:    "newshort12",
			action: func(s *Storage, original, short string) (string, error) {
				ent, err := s.SaveOrGet(ctx, entities.RestoreURL(original, short))
				if ent != nil {
					return ent.Short, err
				}
				return "", err
			},
			expectedVal: "googl12345",
			expectedErr: nil,
		},
		{
			name:  "Get non-existent URL",
			short: "notexists1",
			action: func(s *Storage, original, short string) (string, error) {
				ent, err := s.GetByShort(ctx, short)
				if ent != nil {
					return ent.Original, err
				}
				return "", err
			},
			expectedVal: "",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.action(storage, tt.original, tt.short)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedVal, val)
		})
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()
	var wg sync.WaitGroup
	const numRoutines = 1000

	originals := make([]string, numRoutines)
	shorts := make([]string, numRoutines)

	for i := 0; i < numRoutines; i++ {
		originals[i] = fmt.Sprintf("https://example.com/%d", i)
		shorts[i] = fmt.Sprintf("test%06d", i)
	}

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			urlEntity := entities.RestoreURL(originals[i], shorts[i])
			_, err := storage.SaveOrGet(ctx, urlEntity)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			resEntity, err := storage.GetByShort(ctx, shorts[i])
			assert.NoError(t, err)
			assert.NotNil(t, resEntity)
			assert.Equal(t, originals[i], resEntity.Original)
		}(i)
	}
	wg.Wait()
}
