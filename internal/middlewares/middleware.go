// Package middlewares хранит middleware для web-сервера.
package middlewares

import (
	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
)

// Middlewares - структура, через методы которой вызываются middlewares.
type Middlewares struct {
	cfg configs.Config
	a   authenticator.Authenticator
}

// NewMiddlewares - конструктор для Middlewares.
func NewMiddlewares(cfg configs.Config, a authenticator.Authenticator) Middlewares {
	return Middlewares{
		cfg: cfg,
		a:   a,
	}
}
