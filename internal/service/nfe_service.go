package service

import (
	"context"
	repo "nf-safe/internal/repository/nfe"
)

type NfeService struct{
	repo *repo.NfeRepository
}

func NewNfeService(repository *repo.NfeRepository) *NfeService{
	return &NfeService{repository}
}

func (s *NfeService) Register(ctx context.Context, ){
	
}