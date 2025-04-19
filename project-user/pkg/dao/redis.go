package dao

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"project-user/config"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	log.Println("init redis cache")
	rediscClient := redis.NewClient(config.AppConf.InitRedisOptions())
	Rc = &RedisCache{
		rdb: rediscClient,
	}
}
func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	return rc.rdb.Set(ctx, key, value, expire).Err()
}
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil // Key does not exist
	}
	if err != nil {
		return "", err // Some other error occurred
	}
	return val, nil
}
