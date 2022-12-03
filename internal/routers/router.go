package routers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (chi.Router, error) {
	// создаем роутер
	r := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// получаем хандлер
	h, err := handlers.GetHandler()
	if err != nil {
		return nil, err
	}
	// настроим маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.Post)
		r.Get("/{ID}", h.Get)
	})
	return r, nil
}
