package configs

import "os"

const (
	serverAddress = ":8080"
	serverURL     = "http://localhost:8080"
)

func GetServerAddress() string {
	s, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		return s
	}
	return serverAddress
}

func GetServerURL() string {
	s, ok := os.LookupEnv("BASE_URL")
	if ok {
		return s
	}
	return serverURL
}
