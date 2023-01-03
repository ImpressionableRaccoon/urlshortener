package repositories

import (
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	st, _ := NewMemoryStorage()

	url := "testURL"
	var id string

	testUser := uuid.New()

	t.Run("URL not found", func(t *testing.T) {
		r, err := st.Get("test")
		require.NotNil(t, err)
		assert.Equal(t, "", r)
	})

	t.Run("short link", func(t *testing.T) {
		r, err := st.Add(url, testUser)
		require.Nil(t, err)
		id = r
	})

	t.Run("get testURL", func(t *testing.T) {
		r, err := st.Get(id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	t.Run("get testURL from user URLs", func(t *testing.T) {
		r, err := st.GetUserLinks(testUser)
		require.Nil(t, err)
		assert.Contains(t, r, UserLink{
			ID:  id,
			URL: url,
		})
	})
}
