package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"log"
)

type EnvConfig struct {
	PORT          string
	AuthJwtSecret string
	AthTokenExp   int
	RedisHost     string
	RedisPort     string
	DbUrl         string
	GinMode       string
}

var AppConfig *EnvConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	AppConfig = &EnvConfig{
		PORT:          getEnv("PORT"),
		AuthJwtSecret: getEnv("ACCESS_JWT_SECRET"),
		AthTokenExp:   getEnvAsInt("ACCESS_TOKEN_EXPIRE_IN_MINUTES"),
		DbUrl:         getEnv("DB_URL"),
		RedisHost:     getEnv("REDIS_HOST"),
		RedisPort:     getEnv("REDIS_PORT"),
		GinMode:       getEnv("GIN_MODE"),
	}

}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return val
}

func getEnvAsInt(key string) int {
	valStr := getEnv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer: %v", key, err)
	}
	return val
}
