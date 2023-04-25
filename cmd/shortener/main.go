package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	_, n, _ := net.ParseCIDR(cfg.TrustedSubnet)
	if n == nil {
		_, n, err = net.ParseCIDR("127.0.0.1/32")
		if err != nil {
			panic(err)
		}
	}

	h := handlers.NewHandler(s, cfg.EnableHTTPS, cfg.ServerBaseURL, n)
	a := authenticator.New(cfg)
	m := middlewares.NewMiddlewares(cfg, a)
	r := routers.NewRouter(h, m)

	go func() {
		if cfg.PprofServerAddress == "" {
			log.Println("pprof server address is empty, skipping")
			return
		}

		srv := http.Server{
			ReadHeaderTimeout: time.Second,
		}

		var ln net.Listener
		ln, err = net.Listen("tcp", cfg.PprofServerAddress)
		if err != nil {
			log.Printf("pprof listen failed: %v\n", err)
			return
		}

		err = srv.Serve(ln)
		if err != nil {
			log.Printf("pprof server error: %s\n", err)
		}
	}()

	go func() {
		ln, grpcErr := net.Listen("tcp", cfg.GRPCAdress)
		if grpcErr != nil {
			log.Printf("listen grpc port error: %s\n", grpcErr)
			return
		}

		i := interceptors.New(a)
		g := grpc.NewServer(grpc.UnaryInterceptor(i.AuthUnaryInterceptor))
		pb.RegisterShortenerServer(g, shortener.NewGRPCServer(s, cfg.EnableHTTPS, cfg.ServerBaseURL))

		if grpcErr = g.Serve(ln); grpcErr != nil {
			log.Printf("gRPC server error: %s\n", grpcErr)
			return
		}
	}()

	srv := http.Server{
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}

	var ln net.Listener
	if cfg.EnableHTTPS {
		if cfg.ServerBaseURL == "" {
			panic(errors.New("empty HTTPS domain name"))
		}
		ln = autocert.NewListener(cfg.ServerBaseURL)
	} else {
		ln, err = net.Listen("tcp", cfg.ServerAddress)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		<-sigint

		if shutdownErr := srv.Shutdown(context.Background()); shutdownErr != nil {
			log.Printf("error shutdown server: %v", err)
		}

		if closeErr := s.Close(context.Background()); closeErr != nil {
			log.Printf("error close storage: %v", err)
		}

		close(shutdown)
	}()

	err = srv.Serve(ln)
	if !errors.Is(err, http.ErrServerClosed) {
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
