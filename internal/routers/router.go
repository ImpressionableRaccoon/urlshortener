package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/auth"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/gzip"
)

func NewRouter(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(gzip.Request)
	r.Use(auth.UserCookie)

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
				r.Get("/urls", handler.UserURLs)
			})
		})
	})

	return r
}
