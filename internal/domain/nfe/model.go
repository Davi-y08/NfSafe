package nfe

import (
	"nf-safe/internal/domain/company"
	"time"
	"gorm.io/gorm"
)

type Nfe struct {
	gorm.Model
	CompanyID uint `json:"company_id"`
	Company company.Company `gorm:"foreignKey:CompanyID"`
	Number int `json:"number"`
	Value float64 `json:"value"`
	CreationDataNFE time.Time `json:"creation_data_nfe"`
}