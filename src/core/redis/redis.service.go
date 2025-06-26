package redis

import "time"

type redisService struct{}

func NewRedisService() RedisService {
	return &redisService{}
}

func (r *redisService) Set(key string, value interface{}) error {
	return Client.Set(Ctx, key, value, 0).Err()
}

func (r *redisService) SetWithExpiration(key string, value interface{}, expTimeInMinutes int) error {
	expiration := time.Duration(expTimeInMinutes) * time.Minute
	return Client.Set(Ctx, key, value, expiration).Err()
}

func (r *redisService) Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func (r *redisService) Delete(key string) error {
	return Client.Del(Ctx, key).Err()
}
