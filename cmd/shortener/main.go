package main

import (
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"

	"golang.org/x/crypto/acme/autocert"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printInfo()

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

	if cfg.EnableHTTPS {
		if cfg.HTTPSDomain == "" {
			panic(errors.New("empty HTTPS domain name"))
		}
		panic(http.Serve(autocert.NewListener(cfg.HTTPSDomain), r))
	} else {
		panic(http.ListenAndServe(cfg.ServerAddress, r))
	}
}

func printInfo() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
}
