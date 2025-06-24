package auth

import (
	"boiler-platecode/src/apis/user/repository"

	"gorm.io/gorm"
)

func InitAuthController(db *gorm.DB) *AuthController {
	repo:=repository.NewUserRepository(db)
	service:=NewAuthService(repo)
	controller:=NewAuthController(service)
	return controller
}
