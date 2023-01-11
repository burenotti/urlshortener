package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStorage struct {
	redis   *redis.Client
	linkTTL time.Duration
}

func NewRedisStorage(redis *redis.Client, linkTTL time.Duration) RedisStorage {
	return RedisStorage{
		redis:   redis,
		linkTTL: linkTTL,
	}
}

func (r RedisStorage) SaveShortLink(ctx context.Context, linkID string, url string) error {
	res := r.redis.SetNX(ctx, linkID, url, r.linkTTL)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (r RedisStorage) GetSource(ctx context.Context, linkID string) (string, error) {
	res, err := r.redis.Get(ctx, linkID).Result()
	if err != nil {
		return "", ErrNoSourceUrl
	}
	return res, nil
}
