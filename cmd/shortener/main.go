package main

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"log"
	"net/http"
)

const (
	serverAddress = ":8080" // localhost, порт 8080
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
	// создадим роутер
	r, err := routers.NewRouter(handler)
	if err != nil {
		panic(err)
	}
	// запуск сервера
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
