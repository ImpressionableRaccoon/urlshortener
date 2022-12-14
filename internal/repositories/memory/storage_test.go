package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	st, _ := NewStorage()

	url := "testURL"
	var id string

	t.Run("URL not found", func(t *testing.T) {
		r, err := st.Get("test")
		require.NotNil(t, err)
		assert.Equal(t, "", r)
	})

	t.Run("short link", func(t *testing.T) {
		r, err := st.Add(url)
		require.Nil(t, err)
		id = r
	})

	t.Run("get testURL", func(t *testing.T) {
		r, err := st.Get(id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})
}
