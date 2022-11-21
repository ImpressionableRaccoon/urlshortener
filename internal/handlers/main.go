package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/shorturl"
	"io"
	"net/http"
	"path"
)

// RootHandler — обработчик запроса к корню
func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := path.Base(r.URL.Path)
		if id == "/" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		url, err := shorturl.Get(id)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(307)
		w.Write([]byte{})
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		index, err := shorturl.Make(string(b))
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(index))
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
