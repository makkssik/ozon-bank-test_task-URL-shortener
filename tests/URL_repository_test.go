package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/infrastructure/persistence/memory/repositories"

	"github.com/stretchr/testify/require"
)

func TestStorage_ConcurrentAccess(t *testing.T) {
	storage := repositories.NewStorage()
	ctx := context.Background()
	var wg sync.WaitGroup
	const numRoutines = 1000

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			original := fmt.Sprintf("https://example.com/%d", id)
			short := fmt.Sprintf("test%06d", id)

			urlEntity, err := entities.NewURL(original, short)
			require.NoError(t, err)

			_ = storage.Save(ctx, urlEntity)
		}(i)
	}

	wg.Wait()

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			short := fmt.Sprintf("test%06d", id)
			_, err := storage.GetByShort(ctx, short)
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	for i := 0; i < numRoutines; i++ {
		short := fmt.Sprintf("test%06d", i)
		_, err := storage.GetByShort(ctx, short)
		require.NoError(t, err)
	}

}
