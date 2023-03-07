// Package storage хранит интерфейс и конструктор для хранилища.
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

// Storager - интерфейс для хранилища.
type Storager interface {
	Add(
		ctx context.Context, url repositories.URL, userID repositories.User,
	) (id repositories.ID, err error) // Сократить ссылку.
	Get(
		ctx context.Context, id repositories.ID,
	) (url repositories.URL, deleted bool, err error) // Получить оригинальную ссылку по ID.
	GetUserLinks(
		ctx context.Context, user repositories.User,
	) (links []repositories.UserLink, err error) // Получить все ссылки пользователя.
	DeleteUserLinks(
		ctx context.Context, ids []repositories.ID, user repositories.User,
	) error // Удалить ссылки пользователя.
	Pool(ctx context.Context) (ok bool) // Проверить соединение с базой данных.
}

// StoragerType - int для хранения типа хранилища.
type StoragerType int

const (
	MemoryStorage StoragerType = 1 << iota // Хранилище во временной памяти.
	FileStorage                            // Хранилище в текстовом файле.
	PsqlStorage                            // Хранилище в базе данных Postgres.
)

// NewStorager - конструктор для хранилища.
//
// Сам выберет нужный тип, в зависимости от конфигурации сервера:
//  0. PsqlStorage
//  1. FileStorage
//  2. MemoryStorage
func NewStorager(cfg configs.Config) (Storager, error) {
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

func getStoragerType(cfg configs.Config) StoragerType {
	if cfg.DatabaseDSN != "" {
		return PsqlStorage
	}
	if cfg.FileStoragePath != "" {
		return FileStorage
	}
	return MemoryStorage
}
