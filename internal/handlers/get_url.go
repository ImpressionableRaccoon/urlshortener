package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetURL - обработчик, который переадресует короткую ссылку на исходный URL.
func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url, deleted, err := h.st.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if deleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
