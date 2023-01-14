package handlers

import (
	"net/http"
)

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	if h.st.Pool(r.Context()) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
