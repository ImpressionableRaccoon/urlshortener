package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
