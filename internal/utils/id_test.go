package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenRandomID(t *testing.T) {
	got, err := GenRandomID()
	require.Nil(t, err)

	assert.Equal(t, idLength, len(got))

	for _, char := range got {
		assert.Contains(t, allowedCharacters, string(char))
	}
}

func BenchmarkGenRandomID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenRandomID()
		require.Nil(b, err)
	}
}
