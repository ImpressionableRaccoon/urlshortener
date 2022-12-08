package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"io"
	"net/http"
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

	index, err := h.st.Add(requestData.URL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", configs.GetServerURL(), index)

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
	w.Write(responseJSON)
}
