package handlers

import (
	"errors"
	"net/http"
	"nf-safe/internal/domain/company"
	"nf-safe/internal/service"

	"github.com/gin-gonic/gin"
)


type CompanyHandler struct{
	companyService *service.CompanyService
	userService *service.UserService
}

func NewCompanyHandler(cs *service.CompanyService, us *service.UserService) (*CompanyHandler){
	return &CompanyHandler{companyService: cs, userService: us}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context){
	var dto company.CreateCompanyDto

	if err := c.ShouldBindJSON(&dto); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "corpo inválido",
		})
		return
	}

	user_response, err := h.userService.GetUserById(c.Request.Context(), dto.UserID)

	if err != nil{
		if errors.Is(err, service.ErrUserNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"error": "usuário não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ocorreu um erro interno no servidor",
		})
		return
	}

	if user_response == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "usuário não encontrado",
		})
	}

	erro := h.companyService.CreateCompany(c.Request.Context(), *user_response, dto.Cnpj, dto.Name, dto.RazaoSocial, dto.NomeFantasia, dto.Status) 
}