package app

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"log"
	"net/http"
)

func Start() {
	// создаем хранилище для коротких ссылок
	st := storage.NewStorage()
	// создадим роутер
	r := routers.NewRouter(st)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
