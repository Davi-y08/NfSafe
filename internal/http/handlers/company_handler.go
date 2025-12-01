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

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return
	}


	if err := c.ShouldBindJSON(&dto); err != nil{
		c.JSON(400, gin.H{
			"error": "corpo inválido",
		})
		return
	}

	user_response, err := h.userService.GetUserById(c.Request.Context(), userID.(uint))

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
		return
	}

	new_company, erro := h.companyService.CreateCompany(c.Request.Context(), *user_response, dto.Cnpj, dto.Name, dto.RazaoSocial, dto.NomeFantasia, dto.Status) 

	if erro != nil{
		if errors.Is(erro, service.ErrInvalidCNPJ) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error:": "cnpj inválido",
			})
			return
		}

		if errors.Is(erro, service.ErrDatabase){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "erro no banco de dados",
			})
			return
		}

		if errors.Is(erro, service.ErrCompanyExisting){
			c.JSON(http.StatusConflict, gin.H{
				"error": "empresa já cadastrada",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro interno do servidor",
		})
		return
	}

	if new_company == nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "não foi possível criar a empresa",
		})
	}

}