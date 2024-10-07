package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type CodeCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *CodeCache) Get(ctx context.Context, biz string, phone string) (string, error) {
	key := cache.key(biz, phone)
	val, err := cache.client.Get(key).Bytes()
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (cache *CodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	return cache.client.Set(cache.key(biz, phone), code, cache.expiration).Err()
}

func (cache *CodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
