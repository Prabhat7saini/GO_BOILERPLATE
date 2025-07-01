package redis

import (
	"sync"
	"time"
)



type redisService struct{}

var (
	redisServiceInstance RedisService
	serviceOnce          sync.Once
)

// GetRedisService returns the singleton RedisService instance
func GetRedisService() RedisService {
	serviceOnce.Do(func() {
		// Ensure Redis client is initialized
		Init()
		redisServiceInstance = &redisService{}
	})
	return redisServiceInstance
}

// Implementation of RedisService interface

func (r *redisService) Set(key string, value interface{}) error {
	return GetClient().Set(GetContext(), key, value, 0).Err()
}

func (r *redisService) SetWithExpiration(key string, value interface{}, expTimeInMinutes int) error {
	expiration := time.Duration(expTimeInMinutes) * time.Minute
	return GetClient().Set(GetContext(), key, value, expiration).Err()
}

func (r *redisService) Get(key string) (string, error) {
	return GetClient().Get(GetContext(), key).Result()
}

func (r *redisService) Delete(key string) error {
	return GetClient().Del(GetContext(), key).Err()
}
