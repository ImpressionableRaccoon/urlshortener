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

func NewStorager(cfg *configs.Config) (Storager, error) {
	var err error
	switch getStoragerType(cfg) {
	case PsqlStorage:
		return postgres.NewPsqlStorage(cfg.DatabaseDSN)
	case FileStorage:
		var file *os.File
		file, err = os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
		return disk.NewFileStorage(file)
	default:
		return memory.NewMemoryStorage()
	}
}

func getStoragerType(cfg *configs.Config) StoragerType {
	if cfg.DatabaseDSN != "" {
		return PsqlStorage
	}
	if cfg.FileStoragePath != "" {
		return FileStorage
	}
	return MemoryStorage
}
