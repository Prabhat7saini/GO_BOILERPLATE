package user

import (
	"boiler-platecode/src/apis/user/repository"
	"gorm.io/gorm"
)

// Expose UserController to main
func InitUserController(db *gorm.DB) *UserController {
	repo := repository.NewUserRepository(db)
	service := NewUserService(repo)
	controller := NewUserController(service)
	return controller
}
