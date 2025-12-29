package database

import (
	"context"
	"fmt"
	"log"

	"github.com/cvudumbarainformatika/backend/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis initializes the Redis client
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Connected to Redis")
	return rdb, nil
}
