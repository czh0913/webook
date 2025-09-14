package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr, // 从配置里取地址
	})
	return rdb
}
