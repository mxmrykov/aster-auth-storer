package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
)

type IRedisDc interface {
	Set(ctx context.Context, asid, iaid string) error
	GetIAID(ctx context.Context, login string) (string, error)
	GetAsidUser(ctx context.Context, asid string) (string, error)
}

type RedisDc struct {
	Client          *redis.Client
	MaxPoolInterval time.Duration
	AsidExp         time.Duration
}

func NewRedisDc(cfg *config.DcRedis, user, password string) IRedisDc {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: user,
		Password: password,
		DB:       1,
	})

	return &RedisDc{
		Client:          rdb,
		MaxPoolInterval: cfg.MaxPoolInterval,
		AsidExp:         cfg.AsidExp,
	}
}
