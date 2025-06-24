package auth

import (
	"boiler-platecode/src/apis/auth/domain"
	"boiler-platecode/src/apis/user/repository"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/const/exception"
	"boiler-platecode/src/common/utils"
	"boiler-platecode/src/core/config"
	"context"
	"log"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, reqBody *domain.Login) common.ServiceOutput[*domain.LoginResponse]
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (a *authService) Login(ctx context.Context, reqBody *domain.Login) common.ServiceOutput[*domain.LoginResponse] {

	existingUser, err := a.userRepo.FindByFields(ctx, map[string]interface{}{"email": reqBody.Email}, "id", "name", "is_active", "password")

	if err != nil {
		log.Printf("existing user error%+v", err)
		return common.ServiceOutput[*domain.LoginResponse]{
			Exception: exception.GetException(exception.INTERNAL_SERVER_ERROR),
		}
	}

	if existingUser == nil {
		log.Printf("existing user not found")
		return common.ServiceOutput[*domain.LoginResponse]{
			Exception: exception.GetException(exception.USER_NOT_FOUND),
		}
	}

	checkPassword := utils.CheckPassword(existingUser.Password, reqBody.Password)
	if !checkPassword {
		log.Printf("password not match")
		return common.ServiceOutput[*domain.LoginResponse]{
			Exception: exception.GetException(exception.INVALID_CREDENTIALS),
		}
	}

	secret := config.AppConfig.AuthJwtSecret

	authTokenExp := config.AppConfig.AthTokenExp
	
	token, err := utils.GenerateJwtToken(common.AccessToken, existingUser.ID, authTokenExp, secret)

	if err != nil {
		log.Printf("Generate AuthToken error: %v", err)

		return common.ServiceOutput[*domain.LoginResponse]{
			Exception: exception.GetException(exception.INTERNAL_SERVER_ERROR),
		}
	}
	now := time.Now()
	updateUserErr := a.userRepo.UpdateFields(ctx, map[string]interface{}{"id": existingUser.ID}, map[string]interface{}{"last_login": now})
	if updateUserErr != nil {
		log.Println("Failed to update last login:", err)
		// Optionally return internal error
		return common.ServiceOutput[*domain.LoginResponse]{
			Exception: exception.GetException(exception.INTERNAL_SERVER_ERROR),
		}
	}
	response := &domain.LoginResponse{
		AccessToken: token,
	}
	return common.ServiceOutput[*domain.LoginResponse]{
		Message:        "Login Success",
		OutputData:     response,
		HttpStatusCode: http.StatusOK,
	}
}
