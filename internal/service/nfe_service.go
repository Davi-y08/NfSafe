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
	ErrInvalidId = errors.New("id inválida")
	ErrInDatabaseNFE = errors.New("ocorreu um erro com o banco de dados")
	ErrInCnpjValidation = errors.New("cnpj inválido")
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

func (s *NfeService) GetById(ctx context.Context, id uint) (*nfe.Nfe, error){
	if id == 0{
		return nil, ErrInvalidId
	}

	n, err := s.repo.GetById(ctx, id)

	if err != nil{
		return nil, ErrInDatabaseNFE
	}

	return n, nil
}

func (s *NfeService) GetNFEsByCnpj(ctx context.Context, cnpj string) ([]nfe.Nfe, error){
	if len(cnpj) != 14{
		return nil, ErrInCnpjValidation
	}

	nfes, err := s.repo.GetNFEsByCnpj(ctx, cnpj)

	if err != nil{
		return nil, ErrInDatabaseNFE
	}

	if nfes == nil{
		return nil, nil
	}

	return nfes, nil
}

func (s *NfeService) GetAllNFEs(ctx context.Context) ([]nfe.Nfe, error){
	nfes, err := s.repo.GetAllNFEs(ctx)

	if err != nil{
		return nil, ErrInDatabaseNFE
	}

	if nfes == nil{
		return nil, nil
	}

	return nfes, nil
}