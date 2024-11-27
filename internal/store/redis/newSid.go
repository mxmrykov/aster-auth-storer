package redis

import (
	"context"
)

func (r *RedisDc) Set(ctx context.Context, asid, iaid string) error {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	return r.Client.Set(ctx, asid, iaid, r.AsidExp).Err()
}

func (r *RedisDc) GetAsidUser(ctx context.Context, asid string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	return r.Client.Get(ctx, asid).Result()
}

func (r *RedisDc) GetIAID(ctx context.Context, login string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	res, err := r.Client.Get(ctx, login).Result()

	return res, err
}
