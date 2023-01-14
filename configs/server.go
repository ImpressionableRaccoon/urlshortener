package configs

import (
	"flag"
	"os"
)

var (
	ServerAddress   = ":8080"
	ServerBaseURL   = "http://localhost:8080"
	FileStoragePath = ""
	DatabaseDSN     = ""
	CookieKey       = []byte{14, 180, 4, 236, 208, 28, 133, 5, 116, 159, 137, 123, 80, 176, 209, 179}
)

func Load() {
	loadEnv()
	loadArgs()
}

func loadEnv() {
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		ServerAddress = s
	}

	if s, ok := os.LookupEnv("BASE_URL"); ok {
		ServerBaseURL = s
	}

	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		FileStoragePath = s
	}

	if s, ok := os.LookupEnv("DATABASE_DSN"); ok {
		DatabaseDSN = s
	}
}

func loadArgs() {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "server address")
	flag.StringVar(&ServerBaseURL, "b", ServerBaseURL, "server base url")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "disk storage path")
	flag.StringVar(&DatabaseDSN, "d", DatabaseDSN, "database data source name")

	flag.Parse()
}
