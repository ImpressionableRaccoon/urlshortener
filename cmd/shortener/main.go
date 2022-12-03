package main

import (
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
	st := storage.NewStorage()
	// создадим роутер
	r := routers.NewRouter(st)
	// запуск сервера
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
