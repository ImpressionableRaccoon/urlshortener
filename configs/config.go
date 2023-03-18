// Package configs хранит конфигурацию сервера.
package configs

import (
	"flag"
	"os"
)

// Config - структура для хранения конфигурации сервера.
type Config struct {
	ServerAddress      string // Адрес сервера, по умолчанию ":8080".
	PprofServerAddress string // Адрес сервера профилирования.
	ServerBaseURL      string // URL сервера, по умолчанию "http://localhost:8080".
	FileStoragePath    string // Путь для файлового хранилища.
	DatabaseDSN        string // Адрес базы данных.
	CookieKey          []byte // Ключ для подписи cookie.
}

// NewConfig - конструктор для Config, сам получит и запишет значения.
//
// Приоритет (меньше - приоритетнее):
//  0. аргументы командной строки
//  1. env-переменные
//  2. константы из исходника
func NewConfig() Config {
	cfg := Config{
		ServerAddress: ":8080",
		ServerBaseURL: "http://localhost:8080",
		CookieKey:     []byte{14, 180, 4, 236, 208, 28, 133, 5, 116, 159, 137, 123, 80, 176, 209, 179},
	}

	loadEnv(&cfg)
	loadArgs(&cfg)

	return cfg
}

func loadEnv(cfg *Config) {
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		cfg.ServerAddress = s
	}

	if s, ok := os.LookupEnv("BASE_URL"); ok {
		cfg.ServerBaseURL = s
	}

	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = s
	}

	if s, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = s
	}
}

func loadArgs(cfg *Config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.ServerBaseURL, "b", cfg.ServerBaseURL, "server base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data source name")

	flag.Parse()
}
