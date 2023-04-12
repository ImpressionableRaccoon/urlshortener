package disk

import (
	"context"
	"os"
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

// TestFileStorage - тестируем FileStorage.
func TestFileStorage(t *testing.T) {
	var err error

	filename := "testingStorage"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o777)
	require.NoError(t, err)

	st, err := NewFileStorage(file)
	require.NoError(t, err)

	links := []TestLink{
		{URL: "https://google.com", Delete: true},
		{URL: "https://yandex.ru", Delete: false},
		{URL: "https://practicum.yandex.ru/go-advanced/", Delete: false},
	}

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		var url repositories.URL
		url, _, err = st.Get(context.Background(), "test")
		require.Error(t, err)
		assert.Equal(t, "", url)
	})

	t.Run("short links", func(t *testing.T) {
		var id repositories.ID
		for index, link := range links {
			id, err = st.Add(context.Background(), link.URL, testUser)
			require.NoError(t, err)
			link.ID = id
			links[index] = link
		}
	})

	t.Run("get testURLs", func(t *testing.T) {
		var url repositories.URL
		var deleted repositories.Deleted
		for _, link := range links {
			url, deleted, err = st.Get(context.Background(), link.ID)
			require.NoError(t, err)
			assert.Equal(t, link.URL, url)
			assert.Equal(t, false, deleted)
		}
	})

	t.Run("get testURLs from user URLs", func(t *testing.T) {
		var r []repositories.LinkData
		r, err = st.GetUserLinks(context.Background(), testUser)
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
		err = st.DeleteUserLinks(context.Background(), linksIDs, testUser)
		require.NoError(t, err)
	})

	t.Run("check if only needed URL deleted", func(t *testing.T) {
		var url repositories.URL
		var deleted repositories.Deleted
		for _, link := range links {
			url, deleted, err = st.Get(context.Background(), link.ID)
			require.NoError(t, err)
			assert.Equal(t, link.URL, url)
			assert.Equal(t, link.Delete, deleted)
		}
	})

	err = st.Close(context.Background())
	require.NoError(t, err)

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o777)
	require.NoError(t, err)

	st, err = NewFileStorage(file)
	require.NoError(t, err)

	t.Run("get URLs after restart", func(t *testing.T) {
		var url repositories.URL
		var deleted repositories.Deleted
		for _, link := range links {
			url, deleted, err = st.Get(context.Background(), link.ID)
			require.NoError(t, err)
			assert.Equal(t, link.URL, url)
			assert.Equal(t, link.Delete, deleted)
		}
	})

	err = st.Close(context.Background())
	require.NoError(t, err)

	err = os.Remove(filename)
	require.NoError(t, err)

	t.Run("empty file storage", func(t *testing.T) {
		_, err = NewFileStorage(nil)
		assert.Error(t, err)
	})
}
