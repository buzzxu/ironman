package ironman

import (
	"github.com/buzzxu/ironman/logger"
	"time"

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
		var maxRetries = 5
		var minIdleConns = 100
		if len(conf.ServerConf.Redis.Password) > 0 && conf.ServerConf.Redis.Password != "none" {
			password = conf.ServerConf.Redis.Password
		}
		if conf.ServerConf.Redis.PoolSize > 0 {
			poolSize = conf.ServerConf.Redis.PoolSize
		}
		if conf.ServerConf.Redis.MaxRetries > 0 {
			maxRetries = conf.ServerConf.Redis.MaxRetries
		}
		if conf.ServerConf.Redis.MinIdleConns > 0 {
			minIdleConns = conf.ServerConf.Redis.MinIdleConns
		}
		Redis = redis.NewClient(&redis.Options{
			Addr:         conf.ServerConf.Redis.Addr,
			Password:     password,                 // no password set
			DB:           conf.ServerConf.Redis.DB, // use default DB
			PoolSize:     poolSize,
			MaxRetries:   maxRetries,
			MinIdleConns: minIdleConns,
			DialTimeout:  1 * time.Second,
			ReadTimeout:  500 * time.Millisecond,
			WriteTimeout: 500 * time.Millisecond,
			IdleTimeout:  60 * time.Second,
		})
		_, err := Redis.Ping().Result()
		if err != nil {
			logger.Fatalf("Redis connect error.%s", err.Error())
		}
		logger.Info("Redis connect success")
		if conf.ServerConf.Redis.Stats {
			var ticker *time.Ticker
			ticker = time.NewTicker(5 * time.Minute)
			go func() {
				for {
					select {
					case <-ticker.C:
						RedisStats()
					}
				}
			}()
		}
	} else {
		logger.Warn("Redis未配置,无需连接")
	}
}

func RedisStats() {
	poolStats := Redis.PoolStats()
	logger.Infof("Redis Stats:[TotalConns:%d,IdleConns:%d,StaleConns:%d,Hits:%d,Misses:%d]",
		poolStats.TotalConns,
		poolStats.IdleConns,
		poolStats.StaleConns,
		poolStats.Hits,
		poolStats.Misses)
}

func RedisClose() {
	if Redis != nil {
		Redis.Close()
	}
}
