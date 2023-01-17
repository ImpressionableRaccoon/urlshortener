package memory

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

func TestMemoryStorage(t *testing.T) {
	st, _ := NewMemoryStorage()

	url := "testURL"
	var id string

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		r, _, err := st.Get(context.Background(), "test")
		require.NotNil(t, err)
		assert.Equal(t, "", r)
	})

	t.Run("short link", func(t *testing.T) {
		r, err := st.Add(context.Background(), url, testUser)
		require.Nil(t, err)
		id = r
	})

	t.Run("get testURL", func(t *testing.T) {
		r, _, err := st.Get(context.Background(), id)
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
}
