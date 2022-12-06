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
	// создаем хранилище для коротких ссылок
	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}
	// создаем хендлер
	handler, err := handlers.NewHandler(s)
	if err != nil {
		panic(err)
	}
	// создадем роутер
	r, err := routers.NewRouter(handler)
	if err != nil {
		panic(err)
	}
	// запуск сервера
	log.Fatal(http.ListenAndServe(configs.ServerAddress, r))
}
