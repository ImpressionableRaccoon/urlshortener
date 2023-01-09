package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/auth"
	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

type correlationID = string

type batchRequest struct {
	CorrelationID correlationID    `json:"correlation_id"`
	OriginalURL   repositories.URL `json:"original_url"`
}

type batchResponse struct {
	CorrelationID correlationID    `json:"correlation_id"`
	ShortURL      repositories.URL `json:"short_url"`
}

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := uuid.Parse(r.Context().Value(auth.UserKey{}).(string))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	requestData := make([]batchRequest, 0)

	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	responseData := make([]batchResponse, 0, len(requestData))

	for _, link := range requestData {
		response := batchResponse{
			CorrelationID: link.CorrelationID,
		}

		id, err := h.st.Add(r.Context(), link.OriginalURL, user)
		if !errors.Is(err, repositories.URLAlreadyExists) && err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		response.ShortURL = fmt.Sprintf("%s/%s", configs.ServerBaseURL, id)

		responseData = append(responseData, response)
	}

	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("ShortenURL write failed: %v", err)
	}
}
