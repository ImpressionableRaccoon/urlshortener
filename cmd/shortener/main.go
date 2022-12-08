package main

import (
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

func main() {
	s := storage.NewStorage()
	h := handlers.NewHandler(s)
	r := routers.NewRouter(h)

	log.Fatal(http.ListenAndServe(configs.GetServerAddress(), r))
}
