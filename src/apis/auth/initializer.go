package auth

import (
	"boiler-platecode/src/apis/user/repository"
	"boiler-platecode/src/core/redis"

	"gorm.io/gorm"
)

func InitAuthController(db *gorm.DB) *AuthController {
	repo := repository.NewUserRepository(db)
	redisService := redis.NewRedisService()
	service := NewAuthService(repo,redisService)
	controller := NewAuthController(service)
	return controller
}
