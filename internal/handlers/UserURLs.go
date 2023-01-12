package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/google/uuid"
)

type UserLink struct {
	ShortURL    repositories.URL `json:"short_url"`
	OriginalURL repositories.URL `json:"original_url"`
}

func (h *Handler) UserURLs(w http.ResponseWriter, r *http.Request) {
	user, err := uuid.Parse(r.Context().Value(utils.MyKey("userID")).(string))
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
			ShortURL:    fmt.Sprintf("%s/%s", configs.ServerBaseURL, link.ID),
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
		w.Header().Set("content-type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, err = w.Write(responseJSON)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
