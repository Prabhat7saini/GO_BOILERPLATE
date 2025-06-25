package repository

import (
	"boiler-platecode/src/apis/user/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByFields(ctx context.Context, conditions map[string]interface{}, selectFields ...string) (*entity.User, error)
	UpdateFields(ctx context.Context, conditions map[string]interface{}, fields map[string]interface{}) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByFields(ctx context.Context, conditions map[string]interface{}, selectFields ...string) (*entity.User, error) {
	var user entity.User

	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Optional select
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	// Apply map conditions safely
	query = query.Where(conditions)

	// Execute query
	err := query.First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err 
		}
		return nil, err // real DB error
	}
	return &user, nil
}
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = nil
	fmt.Printf("User before create: %+v\n", user)

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateFields(ctx context.Context, conditions map[string]interface{}, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where(conditions).
		Updates(fields).Error
}
