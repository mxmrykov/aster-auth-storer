package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	"time"
)

type IRedis interface {
	Set(ctx context.Context, asid, iaid string) error
	IsAlive(ctx context.Context, asid string) (bool, error)
	GetAsidUser(ctx context.Context, asid string) (string, error)

	PutIAID(ctx context.Context, login, iaid string) error
	GetIAID(ctx context.Context, login string) (string, error)
	IsIAIDAlive(ctx context.Context, iaid string) (bool, error)
}

type Redis struct {
	Client          *redis.Client
	MaxPoolInterval time.Duration
	AsidExp         time.Duration
}

func NewRedis(cfg *config.DcRedis, user, password string) IRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: user,
		Password: password,
		DB:       0,
	})

	return &Redis{
		Client:          rdb,
		MaxPoolInterval: cfg.MaxPoolInterval,
		AsidExp:         cfg.AsidExp,
	}
}
