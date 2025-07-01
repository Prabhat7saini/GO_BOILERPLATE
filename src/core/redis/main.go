package redis

import (
	"boiler-platecode/src/common/lib/logger"
	"boiler-platecode/src/core/config"
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
	module = "redisModule"
	once   sync.Once
)

// Init initializes the Redis client once (singleton)
func Init() {
	once.Do(func() {
		address := config.AppConfig.RedisHost + ":" + config.AppConfig.RedisPort

		client = redis.NewClient(&redis.Options{
			Addr: address,
		})

		if _, err := client.Ping(ctx).Result(); err != nil {
			logger.Error(module, "❌ Failed to connect to Redis: %v", err)
		} else {
			logger.Info(module, "✅ Redis connected successfully!", "")
		}
	})
}

// GetClient returns the singleton Redis client
func GetClient() *redis.Client {
	if client == nil {
		Init()
	}
	return client
}

// GetContext returns the shared Redis context
func GetContext() context.Context {
	return ctx
}

// Close closes the Redis client
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
