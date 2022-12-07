package handlers

import (
	"fmt"
	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url, err := h.st.Get(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	index, err := h.st.Add(string(b))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s%s", configs.ServerURL, index)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
