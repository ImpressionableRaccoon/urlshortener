package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// RootGetHandler - обработчик GET-запросов к корню
func RootGetHandler(w http.ResponseWriter, r *http.Request) {
	st, err := storage.GetStorage()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url, err := st.Get(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// RootPostHandler - обработчик POST-запросов к корню
func RootPostHandler(w http.ResponseWriter, r *http.Request) {
	st, err := storage.GetStorage()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(b) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	index, err := st.Add(string(b))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8080/" + index

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
