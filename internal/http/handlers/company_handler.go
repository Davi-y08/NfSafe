package handlers

import (
	"nf-safe/internal/service"
	"github.com/gin-gonic/gin"
)


type CompanyHandler struct{
	companyService *service.CompanyService
}

func NewCompanyHandler(cs *service.CompanyService) *CompanyHandler{
	return &CompanyHandler{companyService: cs}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context){
	
}