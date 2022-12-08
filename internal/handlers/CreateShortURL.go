package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
)

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	index, err := h.st.Add(string(b))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", configs.GetServerURL(), index)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
