package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

// Типы, которые использует ShortenURL.
type (
	// ShortenURLRequest - структура запроса к ShortenURL.
	ShortenURLRequest struct {
		URL string `json:"url"` // Исходный URL.
	}

	// ShortenURLResponse - структура ответа от ShortenURL.
	ShortenURLResponse struct {
		Result string `json:"result"` // Сокращенный URL.
	}
)

// ShortenURL - обработчик для создания короткой ссылки через JSON POST body.
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var requestData ShortenURLRequest
	err = json.Unmarshal(b, &requestData)
	if err != nil || requestData.URL == "" {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := authenticator.GetUser(r.Context())
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	id, err := h.st.Add(r.Context(), requestData.URL, user)
	if errors.Is(err, repositories.ErrURLAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := &ShortenURLResponse{
		Result: h.genShortLink(id),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
