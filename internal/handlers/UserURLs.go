package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares/auth"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"

	"github.com/google/uuid"
)

type userLink struct {
	ShortURL    repositories.URL `json:"short_url"`
	OriginalURL repositories.URL `json:"original_url"`
}

func (h *Handler) UserURLs(w http.ResponseWriter, r *http.Request) {
	user, err := uuid.Parse(r.Context().Value(auth.UserKey{}).(string))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	links := h.st.GetUserLinks(user)

	if len(links) == 0 {
		http.Error(w, "[]", http.StatusNoContent)
		return
	}

	linkURLs := make([]userLink, 0)

	for _, link := range links {
		linkURLs = append(linkURLs, userLink{
			ShortURL:    fmt.Sprintf("%s/%s", configs.ServerBaseURL, link.ID),
			OriginalURL: link.URL,
		})
	}

	m, err := json.Marshal(linkURLs)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write(m)
}
