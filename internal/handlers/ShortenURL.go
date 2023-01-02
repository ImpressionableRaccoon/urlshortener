package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
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
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var requestData ShortenURLRequest

	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := uuid.Parse(r.Context().Value("user").(string))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	id, err := h.st.Add(requestData.URL, user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", configs.ServerBaseURL, id)

	response := &ShortenURLResponse{
		Result: url,
	}

	responseJSON, err := json.Marshal(response)
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
