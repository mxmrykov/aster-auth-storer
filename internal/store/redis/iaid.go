package redis

import (
	"context"
	"errors"
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

	return res, err
}
