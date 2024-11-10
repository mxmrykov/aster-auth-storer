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
	IRedis     redis.IRedisDc
	GrpcServer *grpc_server.GrpcServer
}

func NewService(ctx context.Context, cfg *config.AuthStorer, logger *zerolog.Logger) (IService, error) {
	v, err := vault.NewVault(&cfg.Vault)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing vault client")
	}

	secrets, err := v.GetSecretRepo(ctx, cfg.Vault.RedisSecret.Path)

	if err != nil {
		logger.Fatal().Err(err).Msg("error getting dc-redis credentials")
	}

	dc, ac := redis.NewRedisDc(
		&cfg.DcRedis,
		secrets[cfg.Vault.RedisSecret.DcRedisUserName],
		secrets[cfg.Vault.RedisSecret.DcRedisSecretName],
	), redis.NewRedisAc(
		&cfg.DcRedis,
		secrets[cfg.Vault.RedisSecret.AcRedisUserName],
		secrets[cfg.Vault.RedisSecret.AcRedisSecretName],
	)

	grpcServer, err := grpc_server.NewGrpcServer(dc, ac, v, cfg, logger)

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
