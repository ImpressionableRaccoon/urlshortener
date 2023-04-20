// Package shortener хранит grpc-сервер для шортенера.
package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	pb "github.com/ImpressionableRaccoon/urlshortener/proto"
)

type server struct {
	pb.UnimplementedShortenerServer

	cfg configs.Config
	s   storage.Storager
}

// NewGRPCServer - конструктор сервера шортенера.
func NewGRPCServer(cfg configs.Config, s storage.Storager) *server {
	return &server{
		cfg: cfg,
		s:   s,
	}
}

// Ping - обработчик для проверки связи с хранилищем.
func (s server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if !s.s.Pool(ctx) {
		return nil, status.Error(codes.Internal, "database ping failed")
	}

	return &emptypb.Empty{}, nil
}

// Short - обработчик для создания короткой ссылки.
func (s server) Short(ctx context.Context, l *pb.Link) (*pb.Link, error) {
	user, err := authenticator.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	return s.short(ctx, l, user)
}

// Get - обработчик, который получает полную ссылку из id короткой.
func (s server) Get(ctx context.Context, l *pb.Link) (*pb.Link, error) {
	if len(l.Id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "id length should be greater than 0")
	}
	url, deleted, err := s.s.Get(ctx, l.Id)
	if errors.Is(err, repositories.ErrURLNotFound) {
		return nil, status.Error(codes.NotFound, "url not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	if deleted {
		return nil, status.Error(codes.Unavailable, "link is deleted")
	}

	return &pb.Link{
		Id:       l.Id,
		Url:      url,
		ShortUrl: s.genShortLink(l.Id),
	}, nil
}

// GetLinks - обработчик возвращающий все ссылки принадлежащие текущему пользователю.
func (s server) GetLinks(ctx context.Context, _ *emptypb.Empty) (*pb.Batch, error) {
	user, err := authenticator.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	links, err := s.s.GetUserLinks(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	b := &pb.Batch{}
	for _, link := range links {
		b.Data = append(b.Data, &pb.Link{
			Id:       link.ID,
			Url:      link.URL,
			ShortUrl: s.genShortLink(link.ID),
		})
	}

	return b, nil
}

// BatchShort - обработчик для создания пачки коротких ссылок.
func (s server) BatchShort(ctx context.Context, in *pb.Batch) (*pb.Batch, error) {
	user, err := authenticator.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	res := &pb.Batch{}

	for _, link := range in.Data {
		l, err := s.short(ctx, link, user)
		if err != nil {
			continue
		}
		res.Data = append(res.Data, l)
	}

	return res, nil
}

// Delete - обработчик для удаления ссылок пользователя.
func (s server) Delete(ctx context.Context, b *pb.Batch) (*emptypb.Empty, error) {
	user, err := authenticator.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	ids := make([]repositories.ID, 0)
	for _, link := range b.Data {
		if len(link.Id) == 0 {
			continue
		}
		ids = append(ids, link.Id)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := s.s.DeleteUserLinks(ctx, ids, user)
		if err != nil {
			log.Printf("unable to delete user ids: %v", err)
		}
	}()

	return &emptypb.Empty{}, nil
}

// GetStats - обработчик, который возвращает статистику сервера.
func (s server) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.Statistic, error) {
	stats, err := s.s.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &pb.Statistic{
		Links: stats.URLs,
		Users: stats.Users,
	}, nil
}

func (s server) short(ctx context.Context, l *pb.Link, user uuid.UUID) (*pb.Link, error) {
	if len(l.Url) == 0 {
		return nil, status.Error(codes.InvalidArgument, "link length should be greater than 0")
	}
	id, err := s.s.Add(ctx, l.Url, user)
	if errors.Is(err, repositories.ErrURLAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "link already exists")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &pb.Link{
		Id:            id,
		Url:           l.Url,
		ShortUrl:      s.genShortLink(id),
		CorrelationId: l.CorrelationId,
	}, nil
}

func (s server) genShortLink(id string) string {
	if s.cfg.EnableHTTPS {
		return fmt.Sprintf("https://%s/%s", s.cfg.HTTPSDomain, id)
	}
	return fmt.Sprintf("%s/%s", s.cfg.ServerBaseURL, id)
}
