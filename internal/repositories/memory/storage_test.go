package memory

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

// TestLink - структура, которая содержит тестируемую ссылку.
type TestLink struct {
	URL    repositories.URL // Исходный URL.
	ID     repositories.ID  // ID сокращенной ссылки.
	Delete bool             // Удалена ли ссылка.
}

// TestMemoryStorage - тестируем MemStorage.
func TestMemoryStorage(t *testing.T) {
	st, err := NewMemoryStorage()
	require.NoError(t, err)

	links := []TestLink{
		{URL: "https://google.com", Delete: true},
		{URL: "https://yandex.ru", Delete: false},
		{URL: "https://practicum.yandex.ru/go-advanced/", Delete: false},
	}

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		_, _, err := st.Get(context.Background(), "test")
		require.Error(t, err)
	})

	t.Run("short links", func(t *testing.T) {
		for index, link := range links {
			id, err := st.Add(context.Background(), link.URL, testUser)
			require.NoError(t, err)
			link.ID = id
			links[index] = link
		}
	})

	t.Run("get testURLs", func(t *testing.T) {
		for _, link := range links {
			r, deleted, err := st.Get(context.Background(), link.ID)
			require.NoError(t, err)
			assert.Equal(t, link.URL, r)
			assert.Equal(t, false, deleted)
		}
	})

	t.Run("get testURLs from user URLs", func(t *testing.T) {
		r, err := st.GetUserLinks(context.Background(), testUser)
		require.NoError(t, err)
		for _, link := range links {
			assert.Contains(t, r, repositories.LinkData{ID: link.ID, URL: link.URL, User: testUser, Deleted: false})
		}
	})

	t.Run("delete URLs", func(t *testing.T) {
		linksIDs := make([]repositories.ID, 0, len(links))
		for _, link := range links {
			if link.Delete {
				linksIDs = append(linksIDs, link.ID)
				continue
			}
		}
		err := st.DeleteUserLinks(context.Background(), linksIDs, testUser)
		require.NoError(t, err)
	})

	t.Run("check if only needed URL deleted", func(t *testing.T) {
		for _, link := range links {
			r, deleted, err := st.Get(context.Background(), link.ID)
			require.NoError(t, err)
			assert.Equal(t, link.URL, r)
			assert.Equal(t, link.Delete, deleted)
		}
	})
}
