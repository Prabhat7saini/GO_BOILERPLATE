package redis

import (
	"boiler-platecode/src/core/config"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func Init() {

	redisPort:=config.AppConfig.RedisPort
	redisHost:=config.AppConfig.RedisHost
	address := redisHost + ":" + redisPort

	Client = redis.NewClient(&redis.Options{
		Addr:     address,
	})

	if _, err := Client.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully!..")
}
