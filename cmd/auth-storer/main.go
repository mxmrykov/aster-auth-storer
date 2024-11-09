package main

import (
	"context"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	"github.com/mxmrykov/aster-auth-storer/internal/service"
	"github.com/mxmrykov/aster-auth-storer/pkg/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	erg, ctx := errgroup.WithContext(context.Background())

	log.Info().Timestamp().Msg("starting auth-storer service...")
	log.Info().Timestamp().Msg("initializing config...")

	cfg, logger, err := config.InitConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize config")
	}

	log.Info().Timestamp().Msg("initializing service...")

	s, err := service.NewService(ctx, cfg, logger)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize service")
	}

	log.Info().Timestamp().Msg("starting service...")

	erg.Go(func() error {
		return s.Start()
	})

	if err = erg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("failed to start service")
	}

	<-utils.GracefulShutDown()

	log.Info().Timestamp().Msg("graceful shutdown")

	s.Stop()

	log.Info().Timestamp().Msg("service stopped")
}
