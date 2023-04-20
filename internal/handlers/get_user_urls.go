package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

// UserLink - структура ссылки, принадлежащей пользователю.
type UserLink struct {
	ShortURL    repositories.URL `json:"short_url"`    // Сокращенный URL.
	OriginalURL repositories.URL `json:"original_url"` // Исходный URL.
}

// GetUserURLs - обработчик возвращающий все ссылки принадлежащие текущему пользователю.
func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	user, err := authenticator.GetUser(r.Context())
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	links, err := h.st.GetUserLinks(r.Context(), user)
	if err != nil {
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := make([]UserLink, 0)
	for _, link := range links {
		response = append(response, UserLink{
			ShortURL:    h.genShortLink(link.ID),
			OriginalURL: link.URL,
		})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, err = w.Write(responseJSON)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
