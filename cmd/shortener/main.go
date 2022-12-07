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
	s := storage.NewStorage()
	handler := handlers.NewHandler(s)
	r := routers.NewRouter(handler)

	log.Fatal(http.ListenAndServe(configs.ServerAddress, r))
}
