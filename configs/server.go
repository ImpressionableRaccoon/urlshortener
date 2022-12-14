package configs

import "os"

const (
	serverAddress = ":8080"
	serverURL     = "http://localhost:8080"
)

func GetServerAddress() string {
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		return s
	}
	return serverAddress
}

func GetServerURL() string {
	if s, ok := os.LookupEnv("BASE_URL"); ok {
		return s
	}
	return serverURL
}

func GetFileStoragePath() (string, bool) {
	return os.LookupEnv("FILE_STORAGE_PATH")
}
