package configs

import (
	"flag"
	"os"
)

var (
	ServerAddress   = ":8080"
	ServerBaseURL   = "http://localhost:8080"
	FileStoragePath = ""
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
}

func loadArgs() {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "server address")
	flag.StringVar(&ServerBaseURL, "b", ServerBaseURL, "server base url")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "file storage path")

	flag.Parse()
}
