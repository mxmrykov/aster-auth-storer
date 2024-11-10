package redis

import (
	"context"
)

func (r *RedisDc) Set(ctx context.Context, asid, iaid string) error {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	return r.Client.Set(ctx, asid, iaid, r.AsidExp).Err()
}

func (r *RedisDc) IsAlive(ctx context.Context, asid string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	ttl, err := r.Client.TTL(ctx, asid).Result()

	if err != nil {
		return false, err
	}

	return ttl > 0, nil
}

func (r *RedisDc) GetAsidUser(ctx context.Context, asid string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.MaxPoolInterval)

	defer cancel()

	return r.Client.Get(ctx, asid).Result()
}
