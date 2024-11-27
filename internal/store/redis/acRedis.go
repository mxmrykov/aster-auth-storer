package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mxmrykov/aster-auth-storer/internal/config"
	"golang.org/x/net/context"
)

type IRedisAc interface {
	PutIAID(ctx context.Context, login, iaid string) error
	GetIAID(ctx context.Context, login string) (string, error)
}

type RedisAc struct {
	Client          *redis.Client
	MaxPoolInterval time.Duration
}

func NewRedisAc(cfg *config.DcRedis, user, password string) IRedisAc {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: user,
		Password: password,
		DB:       0,
	})

	return &RedisAc{
		Client:          rdb,
		MaxPoolInterval: cfg.MaxPoolInterval,
	}
}
