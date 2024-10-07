package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"jike/internal/domain"
	"time"
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

var ErrKeyNotExist = redis.Nil

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	return user, json.Unmarshal(val, &user)
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	ujson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.client.Set(cache.key(user.Id), ujson, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
