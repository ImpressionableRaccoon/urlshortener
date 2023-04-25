// Package handlers хранит обработчики для http-запросов пользователя.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

// Handler хранит обработчики для http-запросов пользователя.
type Handler struct {
	st      storage.Storager
	https   bool
	domain  string
	trusted *net.IPNet
}

// NewHandler - конструктор для Handler.
func NewHandler(s storage.Storager, https bool, domain string, trusted *net.IPNet) *Handler {
	h := &Handler{
		st:      s,
		https:   https,
		domain:  domain,
		trusted: trusted,
	}

	return h
}

func (h *Handler) httpJSONError(w http.ResponseWriter, msg string, code int) {
	jsonError, _ := json.Marshal(
		struct {
			Error string `json:"error"`
		}{
			Error: msg,
		},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(jsonError)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) genShortLink(id string) string {
	if h.https {
		return fmt.Sprintf("https://%s/%s", h.domain, id)
	}
	return fmt.Sprintf("%s/%s", h.domain, id)
}
