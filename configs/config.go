// Package configs хранит конфигурацию сервера.
package configs

import (
	"encoding/json"
	"flag"
	"io"
	"log"
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
	EnableHTTPS        bool   // Используем ли HTTPS (на 443 порту).
	ConfigFile         string // JSON-файл, в котором хранится конфигурация.
	TrustedSubnet      string // Доверенная сеть, из которой можно получать статистику сервиса.
	GRPCAdress         string // Адрес сервера grpc.
}

// NewConfig - конструктор для Config, сам получит и запишет значения.
//
// Приоритет (меньше - приоритетнее):
//  0. аргументы командной строки
//  1. env-переменные
//  2. JSON-файл с конфигурацией
//  3. константы из исходника
func NewConfig() Config {
	cfg := Config{
		ServerAddress: ":8080",
		ServerBaseURL: "http://localhost:8080",
		CookieKey:     []byte{14, 180, 4, 236, 208, 28, 133, 5, 116, 159, 137, 123, 80, 176, 209, 179},
		GRPCAdress:    ":3200",
	}

	cfg.loadEnv()
	cfg.loadArgs()
	cfg.loadJSON()

	return cfg
}

func (cfg *Config) loadEnv() {
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

	if _, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		cfg.EnableHTTPS = true
	}

	if s, ok := os.LookupEnv("CONFIG"); ok {
		cfg.ConfigFile = s
	}

	if s, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		cfg.TrustedSubnet = s
	}
}

func (cfg *Config) loadArgs() {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.ServerBaseURL, "b", cfg.ServerBaseURL, "server base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data source name")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "enable https support")
	flag.StringVar(&cfg.ConfigFile, "c", cfg.ConfigFile, "JSON config file")
	flag.StringVar(&cfg.ConfigFile, "config", cfg.ConfigFile, "JSON config file")
	flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "trusted subnet")

	flag.Parse()
}

func (cfg *Config) loadJSON() {
	if cfg.ConfigFile == "" {
		return
	}

	c := struct {
		ServerAddress   string `json:"server_address"`
		BaseURL         string `json:"base_url"`
		FileStoragePath string `json:"file_storage_path"`
		DatabaseDSN     string `json:"database_dsn"`
		EnableHTTPS     bool   `json:"enable_https"`
		TrustedSubnet   string `json:"trusted_subnet"`
	}{}

	f, err := os.Open(cfg.ConfigFile)
	if err != nil {
		log.Printf("unable to open config file: %v", err)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("unable to read config file: %v", err)
		return
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Printf("unable to parse config file: %v", err)
		return
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = c.ServerAddress
	}
	if cfg.ServerBaseURL == "" {
		cfg.ServerBaseURL = c.BaseURL
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = c.FileStoragePath
	}
	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = c.DatabaseDSN
	}
	if !cfg.EnableHTTPS {
		cfg.EnableHTTPS = c.EnableHTTPS
	}
	if cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = c.TrustedSubnet
	}
}
