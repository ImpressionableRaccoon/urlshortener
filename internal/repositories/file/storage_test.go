package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	filename := "testingStorage"
	st, err := NewStorage(filename)
	require.Nil(t, err)

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

	t.Run("get test URL", func(t *testing.T) {
		r, err := st.Get(id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	err = st.Close()
	require.Nil(t, err)

	st, err = NewStorage(filename)
	require.Nil(t, err)

	t.Run("get test URL after restart", func(t *testing.T) {
		r, err := st.Get(id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	err = st.Close()
	require.Nil(t, err)

	err = os.Remove(filename)
	require.Nil(t, err)

	t.Run("empty file storage", func(t *testing.T) {
		st, err := NewStorage("")
		assert.NotNil(t, err)
		assert.Nil(t, st)
	})
}
