package service

import (
	"context"
	"errors"
	"nf-safe/internal/domain/company"
	"nf-safe/internal/domain/nfe"
	repo "nf-safe/internal/repository/nfe"
)

type NfeService struct{
	repo *repo.NfeRepository
}

var (
	ErrInvalidNumberNFE = errors.New("numéro inválido")
	ErrInvalidValueNFE = errors.New("valor de nota inválido")
)

func NewNfeService(repository *repo.NfeRepository) *NfeService{
	return &NfeService{repository}
}

func (s *NfeService) Register(ctx context.Context, number int, value float64, c company.Company) error {
	if number == 0{
		return ErrInvalidNumberNFE
	}

	if value < 0.0{
		return ErrInvalidValueNFE
	}

	new_nfe := &nfe.Nfe{
		Number: number,
		CompanyID: c.ID,
		Company: c,
		Value: value,
	}

	return s.repo.Create(ctx, new_nfe)
}