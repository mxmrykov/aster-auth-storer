package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

var ErrorNotFound = errors.New("not found")

func (r *RedisAc) PutIAID(ctx context.Context, login, iaid string) error {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	return r.Client.Set(ctx, login, iaid, 0).Err()
}

func (r *RedisAc) GetIAID(ctx context.Context, login string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	res, err := r.Client.Get(ctx, login).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrorNotFound
		}
		return "", err
	}

	return res, nil
}

func (r *RedisAc) IsIAIDAlive(ctx context.Context, iaid string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	ttl, err := r.Client.TTL(ctx, iaid).Result()

	if err != nil {
		return false, err
	}

	return ttl > 0, nil
}
