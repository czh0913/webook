package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrRedisNotFound = redis.Nil

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: 15 * time.Minute,
	}
}

func (cache *UserCache) Get(ctx context.Context, Id int64) (domain.User, error) {
	key := cache.Key(Id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{}
	err = json.Unmarshal(val, &u)

	return u, nil
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
