package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

type Handler struct {
	st storage.Storager
}

func NewHandler(s storage.Storager) *Handler {
	h := &Handler{
		st: s,
	}

	return h
}
