package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/disk"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/memory"
)

func TestNewStorager(t *testing.T) {
	t.Run("memory storage", func(t *testing.T) {
		got, err := NewStorager(configs.Config{})
		require.NoError(t, err)
		switch got.(type) {
		case *memory.MemStorage:
		default:
			assert.Error(t, errors.New("wrong storager type"))
		}
	})

	t.Run("file storage with wrong file name", func(t *testing.T) {
		_, err := NewStorager(configs.Config{
			FileStoragePath: "/",
		})
		require.Error(t, err)
	})

	t.Run("file storage", func(t *testing.T) {
		fileName := "testStorage"

		got, err := NewStorager(configs.Config{
			FileStoragePath: fileName,
		})
		require.NoError(t, err)
		switch got.(type) {
		case *disk.FileStorage:
		default:
			assert.Error(t, errors.New("wrong storager type"))
		}

		err = os.Remove(fileName)
		require.NoError(t, err)
	})

	t.Run("psql storage with wrong file name", func(t *testing.T) {
		_, err := NewStorager(configs.Config{
			DatabaseDSN: "lalala",
		})
		require.Error(t, err)
	})
}

func Test_getStoragerType(t *testing.T) {
	tests := []struct {
		name string
		cfg  configs.Config
		want StoragerType
	}{
		{
			name: "memory storage",
			cfg:  configs.Config{},
			want: MemoryStorage,
		},
		{
			name: "file storage",
			cfg: configs.Config{
				FileStoragePath: "test",
			},
			want: FileStorage,
		},
		{
			name: "postgres storage",
			cfg: configs.Config{
				DatabaseDSN:     "test",
				FileStoragePath: "test",
			},
			want: PsqlStorage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := getStoragerType(tt.cfg)
			assert.Equal(t, tt.want, res)
		})
	}
}
