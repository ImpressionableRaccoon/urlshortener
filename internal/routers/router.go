// Package routers хранит chi-роутеры для запросов пользователя.
package routers

import (
	"compress/flate"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares"
)

// NewRouter создает новый роутер с нужными обработчиками и middleware.
func NewRouter(handler *handlers.Handler, m middlewares.Middlewares) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))
	r.Use(m.GzipRequest)
	r.Use(m.UserCookie)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.CreateShortURL)
		r.Get("/{ID}", handler.GetURL)

		r.Get("/ping", handler.PingDB)

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", handler.ShortenURL)
				r.Post("/batch", handler.ShortenBatch)
			})

			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", handler.GetUserURLs)
				r.Delete("/urls", handler.DeleteUserURLs)
			})

			r.Route("/internal", func(r chi.Router) {
				r.Get("/stats", handler.GetStats)
			})
		})
	})

	return r
}
