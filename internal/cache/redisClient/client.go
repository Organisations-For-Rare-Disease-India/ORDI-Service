package redisClient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

type RedisConfig struct {
	ADDR     string
	Password string
}

func NewDefaultRedisClient() *redisClient {
	redisConfig := RedisConfig{
		ADDR:     "localhost:6379",
		Password: "", // No password by default
	}
	return NewRedisClient(redisConfig)
}

func NewRedisClient(config RedisConfig) *redisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ADDR,
		Password: config.Password,
		DB:       0,
	})

	return &redisClient{client: rdb}
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
