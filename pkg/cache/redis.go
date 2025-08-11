package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache() *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	return &Cache{rdb: rdb}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	return val, err
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *Cache) Ping() (string, error) {
	res, err := c.rdb.Ping(context.Background()).Result()
	if err != nil {
		return "", fmt.Errorf("failed to connection to redis: %s", err.Error())
	}
	return res, nil
}
