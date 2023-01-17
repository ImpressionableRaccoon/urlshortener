package storage

import (
	"context"
	"os"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/disk"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/memory"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/postgres"
)

type Storager interface {
	Add(ctx context.Context, url repositories.URL, userID repositories.User) (id repositories.ID, err error)
	Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error)
	GetUserLinks(ctx context.Context, user repositories.User) (links []repositories.UserLink, err error)
	DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error
	Pool(ctx context.Context) bool
}

type StoragerType int

const (
	MemoryStorage StoragerType = 1 << iota
	FileStorage
	PsqlStorage
)

func NewStorager() (Storager, error) {
	configs.Load()

	var s Storager
	var err error

	switch getStoragerType() {
	case PsqlStorage:
		s, err = postgres.NewPsqlStorage(configs.DatabaseDSN)
		if err != nil {
			return nil, err
		}
	case FileStorage:
		var file *os.File
		if file, err = os.OpenFile(configs.FileStoragePath, os.O_RDWR|os.O_CREATE, 0777); err != nil {
			return nil, err
		}
		if s, err = disk.NewFileStorage(file); err != nil {
			return nil, err
		}
	default:
		if s, err = memory.NewMemoryStorage(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func getStoragerType() StoragerType {
	if dsn := configs.DatabaseDSN; dsn != "" {
		return PsqlStorage
	} else if path := configs.FileStoragePath; path != "" {
		return FileStorage
	} else {
		return MemoryStorage
	}
}
