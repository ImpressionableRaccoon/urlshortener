package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
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
