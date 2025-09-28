package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"), // 从配置里取地址
	})
	return rdb
}
