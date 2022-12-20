package storage

import (
	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

type Storager interface {
	Add(url string) (id string, err error)
	Get(id string) (string, error)
}

func NewStorager() (Storager, error) {
	var s Storager
	var err error

	configs.Load()

	if path := configs.FileStoragePath; path != "" {
		s, err = repositories.NewFileStorage(path)
	} else {
		s, err = repositories.NewMemoryStorage()
	}
	if err != nil {
		panic(err)
	}

	return s, err
}
