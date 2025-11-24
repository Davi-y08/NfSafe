package nfe

import (
	"nf-safe/internal/domain/company"

	"gorm.io/gorm"
)

type Nfe struct {
	gorm.Model
	CompanyID uint `json:"company_id"`
	Company company.Company `gorm:"foreignKey:CompanyID"`
	Number int `json:"number"`
	Value float64 `json:"value"`
}