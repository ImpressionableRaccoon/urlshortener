package storage

import (
	"os"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

type Storager interface {
	Add(url repositories.URL, userID repositories.User) (repositories.ID, error)
	Get(id repositories.ID) (repositories.URL, error)
	GetUserLinks(user repositories.User) []repositories.UserLink
	IsUserExists(userID repositories.User) bool
}

func NewStorager() (Storager, error) {
	var s Storager
	var err error

	configs.Load()

	if path := configs.FileStoragePath; path != "" {
		var file *os.File

		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}

		s, err = repositories.NewFileStorage(file)
	} else {
		s, err = repositories.NewMemoryStorage()
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}
