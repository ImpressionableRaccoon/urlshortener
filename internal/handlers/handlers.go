// Package handlers хранит обработчики для http-запросов пользователя.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

// Handler хранит обработчики для http-запросов пользователя.
type Handler struct {
	st  storage.Storager
	cfg configs.Config
}

// NewHandler - конструктор для Handler.
func NewHandler(s storage.Storager, cfg configs.Config) *Handler {
	h := &Handler{
		st:  s,
		cfg: cfg,
	}

	return h
}

func (h *Handler) httpJSONError(w http.ResponseWriter, error string, code int) {
	jsonError, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{error})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(jsonError)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) genShortLink(id string) string {
	if h.cfg.EnableHTTPS {
		return fmt.Sprintf("https://%s/%s", h.cfg.HTTPSDomain, id)
	}
	return fmt.Sprintf("%s/%s", h.cfg.ServerBaseURL, id)
}
