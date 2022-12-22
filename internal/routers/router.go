package routers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/gzip"

	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(gzip.Request)
	r.Use(gzip.Response)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.CreateShortURL)
		r.Get("/{ID}", handler.GetURL)

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", handler.ShortenURL)
		})
	})

	return r
}
