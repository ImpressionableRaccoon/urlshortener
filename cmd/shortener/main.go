package main

import (
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
)

func main() {
	s, err := storage.NewStorager()
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler(s)
	r := routers.NewRouter(h)

	log.Fatal(http.ListenAndServe(configs.ServerAddress, r))
}
