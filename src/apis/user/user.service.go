package user

import (
	"boiler-platecode/src/apis/user/domain"
	"boiler-platecode/src/apis/user/entity"
	"boiler-platecode/src/apis/user/repository"
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/const/exception"
	"boiler-platecode/src/common/lib/logger"
	utils "boiler-platecode/src/common/utils"
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"
)


const (module="User Service")
type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) common.ServiceOutput[*struct{}]
	GetUserProfile(ctx context.Context, userId int) common.ServiceOutput[int]
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) common.ServiceOutput[*struct{}] {
	logger.Info(module,"🔍 Checking if user exists:", user.Email)

	// Check for existing user
	existingUser, err := s.repo.FindByFields(ctx, map[string]interface{}{"email": user.Email}, "id", "name")

	if err == nil {
	
		logger.Error(module,"CreateUser:User already exists",errors.New( existingUser.Email))
		return utils.ServiceError[*struct{}](exception.USER_ALREADY_EXISTS)
	}

	if err != gorm.ErrRecordNotFound {
		// Unexpected error
		logger.Error(module,"CreateUser", err,"Error while checking existing user")
		return utils.ServiceError[*struct{}](exception.INTERNAL_SERVER_ERROR)
	}

	hashedPassword := utils.HashPassword(user.Password)
	dbUser := &entity.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		UpdatedAt: nil,
	}

	// Create new user using repository
	if err := s.repo.Create(ctx, dbUser); err != nil {
		return utils.ServiceError[*struct{}](exception.INTERNAL_SERVER_ERROR)
	}

	return common.ServiceOutput[*struct{}]{
		Message:        common.USER_REGISTER_SUCCESS,
		OutputData:     nil,
		HttpStatusCode: http.StatusCreated,
	}
}

func (s *userService) GetUserProfile(ctx context.Context, userId int) common.ServiceOutput[int] {

	return common.ServiceOutput[int]{
		Message:        common.USER_PROFILE,
		OutputData:     userId,
		HttpStatusCode: http.StatusOK,
	}

}
