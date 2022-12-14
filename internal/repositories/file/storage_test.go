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

	err = st.Close()
	require.Nil(t, err)

	st, err = NewStorage(filename)
	require.Nil(t, err)

	t.Run("get testURL after restart", func(t *testing.T) {
		r, err := st.Get(id)
		require.Nil(t, err)
		assert.Equal(t, url, r)
	})

	err = st.Close()
	require.Nil(t, err)

	err = os.Remove(filename)
	require.Nil(t, err)
}
