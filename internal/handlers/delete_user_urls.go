package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

func (h *Handler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user uuid: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	ids := make([]repositories.ID, 0)
	err = json.Unmarshal(b, &ids)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	go func() {
		err = h.st.DeleteUserLinks(context.Background(), ids, user)
		if err != nil {
			log.Printf("unable to delete user ids: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}
