package redis

import (
	"context"
	"ml-in-kube-apiserver/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       int(cfg.DB),
		Username: cfg.User,
	})
	return rdb, rdb.Ping(context.Background()).Err()
}
