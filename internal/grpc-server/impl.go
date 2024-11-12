package grpc_server

import (
	"context"
	"errors"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	ast "github.com/mxmrykov/aster-auth-storer/internal/proto/gen"
	"github.com/mxmrykov/aster-auth-storer/internal/store/redis"
	"github.com/mxmrykov/aster-auth-storer/pkg/clients/vault"
	"github.com/mxmrykov/aster-auth-storer/pkg/sid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	ast.UnimplementedAstServer
	IRedisDc redis.IRedisDc
	IRedisAc redis.IRedisAc
	IVault   vault.IVault
	Cfg      *config.AuthStorer
	Logger   *zerolog.Logger
}

func (s *server) GetIAID(ctx context.Context, in *ast.GetIAIDRequest) (*ast.GetIAIDResponse, error) {
	s.Logger.Info().Msg("New call")
	realCc, err := s.IVault.GetSecret(ctx, s.Cfg.Vault.TokenRepo.Path, s.Cfg.Vault.TokenRepo.AstJwtSecretName)

	if err != nil {
		s.Logger.Err(err).Send()
		return nil, status.Error(codes.Internal, "vault aborted: "+err.Error())
	}

	switch {
	case in.ConfirmCode == "":
		err = status.Error(codes.InvalidArgument, "confirm code is required")
		s.Logger.Err(err).Send()
		return nil, err
	case in.ConfirmCode != realCc:
		err = status.Error(codes.InvalidArgument, "confirm code is incorrect")
		s.Logger.Err(err).Send()
		return nil, err
	case in.Login == "":
		err = status.Error(codes.InvalidArgument, "login is required")
		s.Logger.Err(err).Send()
		return nil, err
	}

	iaid, err := s.IRedisAc.GetIAID(ctx, in.Login)
	asid := sid.New(iaid)

	// example of validation sid
	s.Logger.Info().Msgf("Validate new sid result: %v", sid.Validate(asid))

	if err != nil {
		switch {
		case errors.Is(err, redis.ErrorNotFound):
			return &ast.GetIAIDResponse{
				Has:     false,
				IAID:    "",
				ASID:    asid,
				Message: "no such login",
			}, nil
		default:
			s.Logger.Err(err).Send()
			return nil, status.Error(codes.Internal, "redis aborted: "+err.Error())
		}
	}

	if err = s.IRedisDc.Set(ctx, asid, iaid); err != nil {
		s.Logger.Err(err).Send()
		return nil, status.Error(codes.Internal, "redis aborted: "+err.Error())
	}

	return &ast.GetIAIDResponse{
		Has:  true,
		IAID: asid,
		ASID: asid,
	}, nil
}
