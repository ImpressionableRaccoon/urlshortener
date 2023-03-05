package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

type correlationID = string

type (
	BatchRequest struct {
		CorrelationID correlationID    `json:"correlation_id"`
		OriginalURL   repositories.URL `json:"original_url"`
	}

	BatchResponse struct {
		CorrelationID correlationID    `json:"correlation_id"`
		ShortURL      repositories.URL `json:"short_url"`
	}
)

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	requestData := make([]BatchRequest, 0)
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := make([]BatchResponse, 0, len(requestData))
	var id repositories.ID
	for _, link := range requestData {
		id, err = h.st.Add(r.Context(), link.OriginalURL, user)
		if !errors.Is(err, repositories.ErrURLAlreadyExists) && err != nil {
			h.httpJSONError(w, "Server error", http.StatusInternalServerError)
			return
		}

		response = append(response, BatchResponse{
			CorrelationID: link.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.ServerBaseURL, id),
		})
	}

	responseJSON, err := json.Marshal(&response)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
