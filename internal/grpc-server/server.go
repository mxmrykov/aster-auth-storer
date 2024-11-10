package grpc_server

import (
	"fmt"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	ast "github.com/mxmrykov/aster-auth-storer/internal/proto/gen"
	"github.com/mxmrykov/aster-auth-storer/internal/store/redis"
	"github.com/mxmrykov/aster-auth-storer/pkg/clients/vault"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"time"
)

type IGrpcServer interface {
	Serve() error
	GracefulStop()
}

type GrpcServer struct {
	lis         net.Listener
	S           *grpc.Server
	MaxPollTime time.Duration
}

func NewGrpcServer(dc redis.IRedisDc, ac redis.IRedisAc, vault vault.IVault, cfg *config.AuthStorer, l *zerolog.Logger) (*GrpcServer, error) {
	s := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcServer.Port))

	if err != nil {
		return nil, err
	}

	ast.RegisterAstServer(s, &server{
		IRedisDc: dc,
		IRedisAc: ac,
		IVault:   vault,
		Cfg:      cfg,
		Logger:   l,
	})

	return &GrpcServer{
		lis:         lis,
		S:           s,
		MaxPollTime: cfg.GrpcServer.MaxPollTime,
	}, nil
}

func (s *GrpcServer) Serve() error {
	return s.S.Serve(s.lis)
}

func (s *GrpcServer) GracefulStop() {
	s.S.GracefulStop()
}
