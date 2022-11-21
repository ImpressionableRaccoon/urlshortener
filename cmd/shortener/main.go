package main

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.RootHandler)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
