package gzip

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer func(gz *gzip.Reader) {
				err = gz.Close()
				if err != nil {
					log.Printf("GzipRequest gz.Close() failed: %v", err)
				}
			}(gz)
		}

		next.ServeHTTP(w, r)
	})
}