package handlers

import (
	"net/http"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	if repositories.PoolPSQL() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
