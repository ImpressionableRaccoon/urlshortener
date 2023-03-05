package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := configs.NewConfig()

	s, err := storage.NewStorager(cfg)
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler(s, cfg)

	m := middlewares.NewMiddlewares(cfg)

	r := routers.NewRouter(h, m)

	go func() {
		if cfg.PprofServerAddress == "" {
			log.Println("pprof server address is empty, skipping")
			return
		}
		err := http.ListenAndServe(cfg.PprofServerAddress, nil)
		if err != nil {
			log.Printf("pprof server error: %s\n", err)
		}
	}()

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
