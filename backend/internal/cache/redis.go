package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis initializes Redis connection
func InitRedis(addr, password string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
		PoolSize: 20,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

// Get retrieves value from cache and unmarshals into dest
func Get(ctx context.Context, key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set stores value in cache with expiration
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, data, expiration).Err()
}

// Delete removes keys from cache
func Delete(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// Exists checks if key exists in cache
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := RedisClient.Exists(ctx, key).Result()
	return n > 0, err
}

// Close closes Redis connection
func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}