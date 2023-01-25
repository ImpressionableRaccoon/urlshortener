package middlewares

import (
	"github.com/ImpressionableRaccoon/urlshortener/configs"
)

type Middlewares struct {
	cfg *configs.Config
}

func NewMiddlewares(cfg *configs.Config) Middlewares {
	return Middlewares{
		cfg: cfg,
	}
}
