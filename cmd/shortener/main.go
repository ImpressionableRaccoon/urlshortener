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
	_, err := storage.GetStorage()
	if err != nil {
		panic(err)
	}
	// создадим роутер
	r, err := routers.NewRouter()
	if err != nil {
		panic(err)
	}
	// запуск сервера
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
