//go:build k8s

package config

var Config = config{
	// 本地连接
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:3306)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:11479",
	},
}
