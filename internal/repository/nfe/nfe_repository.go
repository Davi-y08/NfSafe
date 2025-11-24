package repository

import (
	"context"
	"nf-safe/internal/domain/nfe"

	"gorm.io/gorm"
)

type NfeRepository struct {
	db *gorm.DB
}

func NewNfeRepository(db *gorm.DB) *NfeRepository{
	return &NfeRepository{db}
}

func (r *NfeRepository) Create(ctx context.Context, n *nfe.Nfe) error{
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *NfeRepository) GetById(ctx context.Context, id uint) (*nfe.Nfe, error){
	var nf nfe.Nfe

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&nf)

	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return nil, nil
		}

		return nil, result.Error
	}

	return &nf, nil
}

func (r *NfeRepository) GetNFEsByCnpj(ctx context.Context, cnpj string) ([]nfe.Nfe, error){
	var nfes []nfe.Nfe
	result := r.db.WithContext(ctx).Where("cnpj = ?", cnpj).Find(&nfes)

	if result.Error != nil{
		return nil, result.Error
	}

	return nfes, nil
}

func (r *NfeRepository) GetAllNFEs(ctx context.Context) ([]nfe.Nfe, error){
	var nfes []nfe.Nfe
	result := r.db.WithContext(ctx).Find(&nfes)

	if result.Error != nil{
		return nil, result.Error
	}

	return nfes, nil
}



//select * from nfes where cnpj = ? 