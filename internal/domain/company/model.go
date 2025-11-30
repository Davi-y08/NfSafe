package company

import (
	"nf-safe/internal/domain/user"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User user.User `gorm:"foreignKey:UserID"`
	Cnpj string `json:"cnpj" gorm:"uniqueIndex;size:14"`
	Name string `json:"name"`
	RazaoSocial string `json:"razao_social"`
	NomeFantasia string `json:"nome_fantasia"`
	Status string `json:"status"`	
}