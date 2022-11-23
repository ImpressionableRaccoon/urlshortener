package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"io"
	"net/http"
	"path"
)

// RootHandler — обработчик запроса к корню
func RootHandler(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			RootGetHandler(w, r, st)
		case http.MethodPost:
			RootPostHandler(w, r, st)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	}
}

// RootGetHandler - обработчик GET-запросов к корню
func RootGetHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	id := path.Base(r.URL.Path)
	if id == "/" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url, err := st.Get(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(307)
}

// RootPostHandler - обработчик POST-запросов к корню
func RootPostHandler(w http.ResponseWriter, r *http.Request, st storage.Storage) {
	b, err := io.ReadAll(r.Body)
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
