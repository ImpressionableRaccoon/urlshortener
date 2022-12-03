package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

var h *Handler

type Handler struct {
	st *storage.Storage
}

func GetHandler() (*Handler, error) {
	if h != nil {
		return h, nil
	}

	s, err := storage.GetStorage()
	if err != nil {
		return nil, err
	}

	h = &Handler{
		st: s,
	}

	return h, nil
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
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8080/" + index

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
