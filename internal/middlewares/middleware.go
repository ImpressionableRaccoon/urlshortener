// Package middlewares хранит middleware для web-сервера.
package middlewares

import (
	"github.com/ImpressionableRaccoon/urlshortener/configs"
)

// Middlewares - структура, через методы которой вызываются middlewares.
type Middlewares struct {
	cfg configs.Config
}

// NewMiddlewares - конструктор для Middlewares.
func NewMiddlewares(cfg configs.Config) Middlewares {
	return Middlewares{
		cfg: cfg,
	}
}
