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

func TestFileStorage(t *testing.T) {
	filename := "testingStorage"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	st, err := NewFileStorage(file)
	require.Nil(t, err)

	url := "testURL"
	var id string

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		r, err := st.Get(context.Background(), "test")
		require.NotNil(t, err)
		assert.Equal(t, "", r)
	})

	t.Run("short link", func(t *testing.T) {
		r, err := st.Add(context.Background(), url, testUser)
		require.Nil(t, err)
		id = r
	})

	t.Run("get test URL", func(t *testing.T) {
		r, err := st.Get(context.Background(), id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	t.Run("get testURL from user URLs", func(t *testing.T) {
		r, err := st.GetUserLinks(context.Background(), testUser)
		require.Nil(t, err)
		assert.Contains(t, r, repositories.UserLink{
			ID:  id,
			URL: url,
		})
	})

	err = st.Close()
	require.Nil(t, err)

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	st, err = NewFileStorage(file)
	require.Nil(t, err)

	t.Run("get test URL after restart", func(t *testing.T) {
		r, err := st.Get(context.Background(), id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	err = st.Close()
	require.Nil(t, err)

	err = os.Remove(filename)
	require.Nil(t, err)

	t.Run("empty disk storage", func(t *testing.T) {
		st, err := NewFileStorage(nil)
		assert.NotNil(t, err)
		assert.Nil(t, st)
	})
}