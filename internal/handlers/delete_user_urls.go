package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

func (h *Handler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	user, err := uuid.Parse(r.Context().Value(utils.ContextKey("userID")).(string))
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

	log.Println(string(b))

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
