package service

import (
	"context"
	"errors"
	"fmt"
	"nf-safe/internal/domain/company"
	"nf-safe/internal/domain/user"
	repo "nf-safe/internal/repository/company"
)

type CompanyService struct {
	repo *repo.CompanyRepository
}

var (
	ErrInvalidCNPJ = errors.New("cnpj inválido")
	ErrInDatabaseCompany = errors.New("ocorreu um erro com o banco de dados")
	ErrCompanyExisting = errors.New("empresa já cadastrada")
	ErrNotFoundCompany = errors.New("empresa não encontrada")
)

func NewCompanyService(repository *repo.CompanyRepository) *CompanyService{
	return &CompanyService{repository}
}

func (s *CompanyService) CreateCompany(ctx context.Context , u user.User, cnpj, name string) error {	
	
	if len(cnpj) != 14{
		return ErrInvalidCNPJ
	}

	existing, erro := s.repo.GetByCnpj(ctx, cnpj)

	if erro != nil{
		fmt.Println("Erro GetByCnpj:", erro)
		return ErrInDatabaseCompany
	}

	if existing != nil{
		return ErrCompanyExisting
	}

	new_company := &company.Company{
		Name: name,
		UserID: u.ID,
		User: u,
		Cnpj: cnpj,
	}

	return s.repo.Create(ctx, new_company)
}

func (s *CompanyService) GetByCnpj(ctx context.Context, cnpj string) (*company.Company, error) {
	if len(cnpj) == 0{
		return nil, ErrInvalidCNPJ
	}

	comp, err := s.repo.GetByCnpj(ctx, cnpj) 

	if err != nil{
		return nil, ErrInDatabaseCompany
	}

	if comp == nil{
		return nil, ErrNotFoundCompany
	}

	return comp, nil
}