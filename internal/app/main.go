package app

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func Start() {
	// создаем хранилище для коротких ссылок
	st := storage.NewStorage()
	// создаем роутер
	r := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// маршрутизация запросов обработчику
	r.HandleFunc("/", handlers.RootHandler(st))
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
