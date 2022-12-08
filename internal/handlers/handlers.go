package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

type Handler struct {
	st *storage.Storage
}

func NewHandler(s *storage.Storage) *Handler {
	h := &Handler{
		st: s,
	}

	return h
}
