package user

import (
	"boiler-platecode/src/apis/user/domain"
	"boiler-platecode/src/apis/user/entity"
	"boiler-platecode/src/apis/user/repository"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/const/exception"
	utils "boiler-platecode/src/common/utils"
	"context"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) common.ServiceOutput[*domain.User]
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}




func (s *userService) CreateUser(ctx context.Context, user *domain.User) common.ServiceOutput[*domain.User] {
	fmt.Println("🔍 Checking if user exists:", user.Email)

	// Check for existing user
	existingUser, err := s.repo.FindByFields(ctx, map[string]interface{}{"email": user.Email}, "id", "name")

	if err == nil {
		// Found user — already exists
		log.Println("User already exists:", existingUser.Email)
		return common.ServiceOutput[*domain.User]{
			Exception: exception.GetException(exception.USER_ALREADY_EXISTS),
		}
	}

	if err != gorm.ErrRecordNotFound {
		// Unexpected error
		log.Println("Error while checking existing user:", err)
		return common.ServiceOutput[*domain.User]{
			Exception: exception.GetException(exception.INTERNAL_SERVER_ERROR),
		}
	}


	
	hashedPassword := utils.HashPassword(user.Password)
	dbUser := &entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		UpdatedAt: nil,
	}

	// Create new user using repository
	if err := s.repo.Create(ctx, dbUser); err != nil {
		return common.ServiceOutput[*domain.User]{
			Exception: exception.GetException(exception.INTERNAL_SERVER_ERROR),
		}
	}

	return common.ServiceOutput[*domain.User]{
		Message:        common.USER_REGISTER_SUCCESS,
		OutputData:     &domain.User{Name: dbUser.Name, Email: dbUser.Email},
		HttpStatusCode: http.StatusCreated,
	}
}
