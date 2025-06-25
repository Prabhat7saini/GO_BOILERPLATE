package auth

import (
	"boiler-platecode/src/apis/auth/domain"
	"boiler-platecode/src/apis/user/repository"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/const/exception"
	"boiler-platecode/src/common/utils"
	"boiler-platecode/src/core/config"
	"boiler-platecode/src/core/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, reqBody *domain.Login) common.ServiceOutput[*domain.LoginResponse]
}

type authService struct {
	userRepo     repository.UserRepository
	redisService redis.RedisService
}

func NewAuthService(userRepo repository.UserRepository, redisService *redis.RedisService) AuthService {
	return &authService{
		userRepo:     userRepo,
		redisService: *redisService,
	}
}
func (a *authService) Login(ctx context.Context, reqBody *domain.Login) common.ServiceOutput[*domain.LoginResponse] {
	existingUser, err := a.userRepo.FindByFields(
		ctx,
		map[string]interface{}{"email": reqBody.Email},
		"id", "name", "is_active", "password",
	)
	if err != nil {
		log.Printf("Error fetching user: %+v", err)
		return utils.ServiceError[*domain.LoginResponse](exception.INTERNAL_SERVER_ERROR)
	}

	if existingUser == nil {
		log.Println("User not found")
		return utils.ServiceError[*domain.LoginResponse](exception.USER_NOT_FOUND)
	}

	if !utils.CheckPassword(existingUser.Password, reqBody.Password) {
		log.Println("Invalid password")
		return utils.ServiceError[*domain.LoginResponse](exception.INVALID_CREDENTIALS)
	}

	authToken, err := utils.GenerateJwtToken(
		common.AccessToken,
		existingUser.ID,
		config.AppConfig.AthTokenExp,
		config.AppConfig.AuthJwtSecret,
	)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return utils.ServiceError[*domain.LoginResponse](exception.INTERNAL_SERVER_ERROR)
	}

	refreshToken, err := utils.GenerateJwtToken(
		common.RefreshToken,
		existingUser.ID,
		config.AppConfig.RefreshTokenExp,
		config.AppConfig.RefreshJwtSecret,
	)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return utils.ServiceError[*domain.LoginResponse](exception.INTERNAL_SERVER_ERROR)
	}

	if err := a.userRepo.UpdateFields(
		ctx,
		map[string]interface{}{"id": existingUser.ID},
		map[string]interface{}{"last_login": time.Now()},
	); err != nil {
		log.Printf("Failed to update last login: %v", err)
		return utils.ServiceError[*domain.LoginResponse](exception.INTERNAL_SERVER_ERROR)
	}

	if err := a.redisService.SetWithExpiration(
		fmt.Sprintf("Auth:userId:%d", existingUser.ID),
		authToken,
		10,
	); err != nil {
		log.Printf("Failed to cache token in Redis: %v", err)
	}
	return common.ServiceOutput[*domain.LoginResponse]{
		Message:        "Login Success",
		OutputData:     &domain.LoginResponse{AccessToken: authToken, RefreshToken: refreshToken},
		HttpStatusCode: http.StatusOK,
	}
}
