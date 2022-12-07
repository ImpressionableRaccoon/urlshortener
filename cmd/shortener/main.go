package main

import (
	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	handler, err := handlers.NewHandler(s)
	if err != nil {
		panic(err)
	}

	r, err := routers.NewRouter(handler)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(configs.ServerAddress, r))
}
