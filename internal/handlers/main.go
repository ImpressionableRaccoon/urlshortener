package handlers

import (
	"fmt"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
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

		url, err := storage.Get(id)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(307)
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		index, err := storage.Make(string(b))
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		url := "http://localhost:8080/" + index

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(url))
		if err != nil {
			fmt.Println(err)
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
