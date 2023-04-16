package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
)

// GetStats - обработчик, который возвращает статистику сервера при запросах из внутренней сети.
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}
	ip := net.ParseIP(addr)
	if h.cfg.TrustedSubnet == nil || !h.cfg.TrustedSubnet.Contains(ip) {
		h.httpJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	stats, err := h.st.GetStats(r.Context())
	if err != nil {
		log.Printf("unable to get stats: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(stats)
	if err != nil {
		log.Printf("unable to marshal response: %v", err)
		h.httpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
