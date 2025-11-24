package repository

import (
	"nf-safe/internal/domain/company"
	"context"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyrepository(db *gorm.DB) *CompanyRepository{
	return &CompanyRepository{db}
}

func (r *CompanyRepository) Create(ctx context.Context, c *company.Company) error{
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CompanyRepository) GetByCnpj(ctx context.Context, cnpj string) (*company.Company, error){
	var comp company.Company 

	result := r.db.WithContext(ctx).Where("cnpj = ?", cnpj).First(&comp)

	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return nil, nil
		}

		return nil, result.Error
	}

	return &comp, nil
}

func (r *CompanyRepository) GetById(ctx context.Context, id uint) (*company.Company, error){
	var comp company.Company 

	result := r.db.WithContext(ctx).First(&comp, id)

	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return nil, nil
		}

		return nil, result.Error
	}

	return &comp, nil
}

func (r *CompanyRepository) GetAllCompanies(ctx context.Context) ([]company.Company, error){
	var companies []company.Company
	result := r.db.WithContext(ctx).Find(&companies)

	if result.Error != nil{
		return nil, result.Error
	}

	return companies, nil
}