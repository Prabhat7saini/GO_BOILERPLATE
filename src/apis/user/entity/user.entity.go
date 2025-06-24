package entity

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	LastLogin *time.Time     `gorm:"default:null" json:"last_login"` // nullable timestamp
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
