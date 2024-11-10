package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mxmrykov/aster-auth-storer/pkg/logger"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type (
	AuthStorer struct {
		UseStackTrace bool       `yaml:"useStackTrace"`
		DcRedis       DcRedis    `yaml:"dcRedis"`
		Vault         Vault      `yaml:"vault"`
		GrpcServer    GrpcServer `yaml:"grpcServer"`
	}

	DcRedis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`

		MaxPoolInterval time.Duration `yaml:"maxPoolInterval"`
		AsidExp         time.Duration `yaml:"asidExp"`
	}

	Vault struct {
		AuthToken     string        `env:"VAULT_AUTH_TOKEN"`
		Host          string        `yaml:"host"`
		Port          int           `yaml:"port"`
		ClientTimeout time.Duration `yaml:"clientTimeout"`

		TokenRepo struct {
			Path string `yaml:"path"`

			AstJwtSecretName string `yaml:"astJwtSecretName"`
		} `yaml:"tokenRepo"`

		RedisSecret struct {
			Path string `yaml:"path"`

			DcRedisSecretName string `yaml:"dcRedisSecretName"`
			DcRedisUserName   string `yaml:"dcRedisUserName"`
			AcRedisSecretName string `yaml:"acRedisSecretName"`
			AcRedisUserName   string `yaml:"acRedisUserName"`
		} `yaml:"dcRedisSecret"`
	}

	GrpcServer struct {
		Port        int           `yaml:"port"`
		MaxPollTime time.Duration `yaml:"maxPollTime"`
	}
)

func InitConfig() (*AuthStorer, *zerolog.Logger, error) {
	cfg := *new(AuthStorer)

	if os.Getenv("BUILD_ENV") == "" {
		return nil, nil, errors.New("build environment is not assigned")
	}

	path := fmt.Sprintf("./deploy/%s.yaml", os.Getenv("BUILD_ENV"))

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, nil, err
	}

	l := logger.NewLogger(cfg.UseStackTrace)

	return &cfg, l, nil
}
