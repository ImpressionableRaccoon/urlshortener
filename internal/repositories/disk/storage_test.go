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

type TestLink struct {
	URL    repositories.URL
	ID     repositories.ID
	Delete bool
}

func TestFileStorage(t *testing.T) {
	filename := "testingStorage"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	st, err := NewFileStorage(file)
	require.Nil(t, err)

	links := []TestLink{
		{URL: "https://google.com", Delete: true},
		{URL: "https://yandex.ru", Delete: false},
		{URL: "https://practicum.yandex.ru/go-advanced/", Delete: false},
	}

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		r, _, err := st.Get(context.Background(), "test")
		require.NotNil(t, err)
		assert.Equal(t, "", r)
	})

	t.Run("short links", func(t *testing.T) {
		for index, link := range links {
			id, err := st.Add(context.Background(), link.URL, testUser)
			require.Nil(t, err)
			link.ID = id
			links[index] = link
		}
	})

	t.Run("get testURLs", func(t *testing.T) {
		for _, link := range links {
			r, deleted, err := st.Get(context.Background(), link.ID)
			require.Nil(t, err)
			assert.Equal(t, link.URL, r)
			assert.Equal(t, false, deleted)
		}
	})

	t.Run("get testURLs from user URLs", func(t *testing.T) {
		r, err := st.GetUserLinks(context.Background(), testUser)
		require.Nil(t, err)
		for _, link := range links {
			assert.Contains(t, r, repositories.UserLink{
				ID:  link.ID,
				URL: link.URL,
			})
		}
	})

	t.Run("delete URLs", func(t *testing.T) {
		linksIDs := make([]repositories.ID, 0)
		for _, link := range links {
			if link.Delete {
				linksIDs = append(linksIDs, link.ID)
				continue
			}
		}
		err := st.DeleteUserLinks(context.Background(), linksIDs, testUser)
		require.Nil(t, err)
	})

	t.Run("check if only needed URL deleted", func(t *testing.T) {
		for _, link := range links {
			r, deleted, err := st.Get(context.Background(), link.ID)
			require.Nil(t, err)
			assert.Equal(t, link.URL, r)
			assert.Equal(t, link.Delete, deleted)
		}
	})

	err = st.Close()
	require.Nil(t, err)

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	st, err = NewFileStorage(file)
	require.Nil(t, err)

	t.Run("get URLs after restart", func(t *testing.T) {
		for _, link := range links {
			r, deleted, err := st.Get(context.Background(), link.ID)
			require.Nil(t, err)
			assert.Equal(t, link.URL, r)
			assert.Equal(t, link.Delete, deleted)
		}
	})

	err = st.Close()
	require.Nil(t, err)

	err = os.Remove(filename)
	require.Nil(t, err)

	t.Run("empty file storage", func(t *testing.T) {
		st, err := NewFileStorage(nil)
		assert.NotNil(t, err)
		assert.Nil(t, st)
	})
}
