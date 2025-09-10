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

/*
A用到了 B，B 一定是接口 =>这个是保证面向接口
A用到了 B，B 一定是 A 的字段 =>规避包变量、包方法，都非常缺乏扩展性
A用到了 B，A 绝对不初始化 B，而是外面注入 =>保持依赖注入(DI，Dependency Injection能据举例子这个是什么意思么我看不懂
*/

type UserCache interface {
	Get(ctx context.Context, Id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: 15 * time.Minute,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, Id int64) (domain.User, error) {
	key := cache.Key(Id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{}
	err = json.Unmarshal(val, &u)

	return u, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
