package handlers

import (
	"fmt"
	"github.com/ImpressionableRaccoon/urlshortener/internal/shorturl"
	"io"
	"net/http"
	"path"
)

// RootHandler — обработчик запроса к корню
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nNew request:", r.Method)
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

		fmt.Println("method GET")
		fmt.Println("Path:", path.Base(r.URL.Path))
		fmt.Println("URL:", url)
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
		_, err = w.Write([]byte(index))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("method POST")
		fmt.Println("Content-Type:", r.Header.Get("Content-Type"))
		fmt.Println("Body:", string(b))
		fmt.Println("Response:", index)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Println("unknown")
	}
}
