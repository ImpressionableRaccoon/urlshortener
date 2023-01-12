package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/google/uuid"
)

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	value := r.Context().Value(utils.MyKey("userID")).(string)
	user, err := uuid.Parse(value)
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var conflict bool
	id, err := h.st.Add(r.Context(), string(b), user)
	if errors.Is(err, repositories.ErrURLAlreadyExists) {
		conflict = true
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", configs.ServerBaseURL, id)

	w.Header().Set("content-type", "text/plain; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if conflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = w.Write([]byte(url))
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
