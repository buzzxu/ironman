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
	if conf.ServerConf.Redis != nil && conf.ServerConf.Redis.Addr != "" {
		var password = ""
		var poolSize = conf.ServerConf.MaxProc * 5
		if len(conf.ServerConf.Redis.Password) > 0 && conf.ServerConf.Redis.Password != "none" {
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
	} else {
		fmt.Printf("Redis未配置,无需连接")
	}
}

func RedisClose() {
	if Redis != nil {
		Redis.Close()
	}
}
