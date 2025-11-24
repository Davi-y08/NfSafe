package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex;size:180"`
	PasswordHash string `json:"-"`
	Role string `json:"role" gorm:"default:user"`
}