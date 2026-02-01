package Cache

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *Config.Config) {
	redisConfig := cfg.Redis
	redisClient = redis.NewClient(&redis.Options{
		ReadTimeout:  redisConfig.ReadTimeOut * time.Second,
		WriteTimeout: redisConfig.WriteTimeOut * time.Second,
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.Db,
		DialTimeout:  redisConfig.DialTimeOut * time.Second,
		PoolSize:     redisConfig.PoolSize,
		PoolTimeout:  redisConfig.PoolTimeOut * time.Second,
	})
}

func GetRedisInstance() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}
