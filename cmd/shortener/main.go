package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/grpc/interceptors"
	"github.com/ImpressionableRaccoon/urlshortener/internal/grpc/shortener"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares"
	"github.com/ImpressionableRaccoon/urlshortener/internal/routers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	pb "github.com/ImpressionableRaccoon/urlshortener/proto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	shutdown := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	printInfo()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := configs.NewConfig()

	s, err := storage.NewStorager(cfg)
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler(s, cfg)
	a := authenticator.New(cfg)
	m := middlewares.NewMiddlewares(cfg, a)
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

	go func() {
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Printf("listen grpc port error: %s\n", err)
			return
		}

		i := interceptors.New(a)
		g := grpc.NewServer(grpc.UnaryInterceptor(i.AuthUnaryInterceptor))
		pb.RegisterShortenerServer(g, shortener.NewGRPCServer(cfg, s))

		if err := g.Serve(listen); err != nil {
			log.Printf("gRPC server error: %s\n", err)
			return
		}
	}()

	srv := http.Server{
		Handler: r,
	}

	var ln net.Listener
	if cfg.EnableHTTPS {
		if cfg.HTTPSDomain == "" {
			panic(errors.New("empty HTTPS domain name"))
		}
		ln = autocert.NewListener(cfg.HTTPSDomain)
	} else {
		ln, err = net.Listen("tcp", cfg.ServerAddress)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		<-sigint

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("error shutdown server: %v", err)
		}

		err = s.Close(context.Background())
		if err != nil {
			log.Printf("error close storage: %v", err)
		}

		close(shutdown)
	}()

	err = srv.Serve(ln)
	if err != http.ErrServerClosed {
		panic(err)
	}

	<-shutdown

	log.Print("shutdown successful")
}

func printInfo() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
}
