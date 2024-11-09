package service

import (
	"context"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	grpc_server "github.com/mxmrykov/aster-auth-storer/internal/grpc-server"
	"github.com/mxmrykov/aster-auth-storer/internal/store/redis"
	"github.com/mxmrykov/aster-auth-storer/pkg/clients/vault"
	"github.com/rs/zerolog"
)

type IService interface {
	Start() error
	Stop()
}

type Service struct {
	Zerolog *zerolog.Logger
	Cfg     *config.AuthStorer

	Vault      vault.IVault
	IRedis     redis.IRedis
	GrpcServer *grpc_server.GrpcServer
}

func NewService(ctx context.Context, cfg *config.AuthStorer, logger *zerolog.Logger) (IService, error) {
	v, err := vault.NewVault(&cfg.Vault)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing vault client")
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing etcd client")
	}

	dcSecrets, err := v.GetSecretRepo(ctx, cfg.Vault.DpRedisSecret.Path)

	if err != nil {
		logger.Fatal().Err(err).Msg("error getting dc-redis password")
	}

	dc := redis.NewRedis(
		&cfg.DcRedis,
		dcSecrets[cfg.Vault.DpRedisSecret.DpRedisUserName],
		dcSecrets[cfg.Vault.DpRedisSecret.DpRedisSecretName],
	)

	grpcServer, err := grpc_server.NewGrpcServer(dc, v, cfg, logger)

	return &Service{
		Zerolog:    logger,
		Cfg:        cfg,
		Vault:      v,
		IRedis:     dc,
		GrpcServer: grpcServer,
	}, nil
}

func (s *Service) Start() error {
	return s.GrpcServer.Serve()
}
func (s *Service) Stop() {
	s.GrpcServer.GracefulStop()
}
