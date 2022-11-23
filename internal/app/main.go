package app

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"log"
	"net/http"
)

func Start() {
	// создаем хранилище для коротких ссылок
	st := storage.NewStorage()
	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.RootHandler(st))
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
