package Cache

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *Config.Config, ctx context.Context) error {
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

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisInstance() *redis.Client {
	return redisClient
}

func Set[T any](redisClient *redis.Client, key string, value T, duration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), key, v, duration).Err()
}

func Get[T any](redisClient *redis.Client, key string) (T, error) {
	var dest T = *new(T)

	v, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return dest, err
	}

	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}

	return dest, nil
}

func CloseRedis() {
	redisClient.Close()
}
