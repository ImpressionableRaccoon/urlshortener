package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	user, err := uuid.Parse(r.Context().Value(utils.ContextKey("userID")).(string))
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	var conflict bool
	id, err := h.st.Add(r.Context(), requestData.URL, user)
	if errors.Is(err, repositories.ErrURLAlreadyExists) {
		conflict = true
	} else if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := &ShortenURLResponse{
		Result: fmt.Sprintf("%s/%s", configs.ServerBaseURL, id),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if conflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
