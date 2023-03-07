package handlers

import (
	"log"
	"net/http"
)

// PingDB - обработчик для проверки связи с хранилищем.
func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	if !h.st.Pool(r.Context()) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
