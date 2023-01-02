package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/auth"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
)

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	value := r.Context().Value(auth.UserKey{}).(string)
	user, err := uuid.Parse(value)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	id, err := h.st.Add(string(b), user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", configs.ServerBaseURL, id)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(url))
	if err != nil {
		log.Printf("CreateShortURL write failed: %v", err)
	}
}
