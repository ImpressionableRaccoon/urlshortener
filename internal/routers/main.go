package routers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(st storage.Storage) chi.Router {
	// создаем роутер
	r := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// настроим маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.RootPostHandler(w, r, st)
		})
		r.Get("/{ID}", func(w http.ResponseWriter, r *http.Request) {
			handlers.RootGetHandler(w, r, st)
		})
	})
	return r
}
