package ironman

import (
	"fmt"

	"github.com/buzzxu/ironman/conf"
	"github.com/go-redis/redis"
)

// Redis 客户端
var Redis *redis.Client

// RedisConnect Redis连接
func RedisConnect() {
	var password = ""
	var poolSize = 10
	if len(conf.ServerConf.Redis.Password) > 0 {
		password = conf.ServerConf.Redis.Password
	}
	if conf.ServerConf.Redis.PoolSize > 0 {
		poolSize = conf.ServerConf.Redis.PoolSize
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     conf.ServerConf.Redis.Addr,
		Password: password,                 // no password set
		DB:       conf.ServerConf.Redis.DB, // use default DB
		PoolSize: poolSize,
	})

	pong, err := Redis.Ping().Result()
	fmt.Println(pong, err)
}
