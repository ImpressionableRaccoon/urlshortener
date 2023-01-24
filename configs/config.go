package configs

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	ServerBaseURL   string
	FileStoragePath string
	DatabaseDSN     string
	CookieKey       []byte
}

func NewConfig() *Config {
	cfg := &Config{
		ServerAddress:   ":8080",
		ServerBaseURL:   "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
		CookieKey:       []byte{14, 180, 4, 236, 208, 28, 133, 5, 116, 159, 137, 123, 80, 176, 209, 179},
	}

	loadEnv(cfg)
	loadArgs(cfg)

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
