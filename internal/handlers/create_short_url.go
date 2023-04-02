package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

// CreateShortURL - обработчик для создания короткой ссылки через обычный POST body.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	id, err := h.st.Add(r.Context(), string(b), user)
	if errors.Is(err, repositories.ErrURLAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := h.genShortLink(id)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(url))
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
