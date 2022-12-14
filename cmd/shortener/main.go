package main

import (
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/memory"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/file"

	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

func main() {
	var s storage.Storage
	var err error

	configs.Load()

	if path := configs.FileStoragePath; path != "" {
		s, err = file.NewStorage(path)
	} else {
		s, err = memory.NewStorage()
	}
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler(s)
	r := routers.NewRouter(h)

	log.Fatal(http.ListenAndServe(configs.ServerAddress, r))
}
