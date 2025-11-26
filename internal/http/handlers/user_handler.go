package handlers

import (
	"errors"
	"net/http"
	"nf-safe/internal/domain/user"
	"nf-safe/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}
////////////////////////////////////////////////////////////////////////////////
// Fazer a implementação de: paginations, auth, async
func NewUserHandler(us *service.UserService) *UserHandler{
	return &UserHandler{userService: us}
}

func (h *UserHandler) CreateUser(c *gin.Context){
	var dto user.CreateUserDto
	
	if err := c.ShouldBindJSON(&dto); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": "corpo inválido",
		})
		return
	}
	
	new_user, err := h.userService.CreateUser(c.Request.Context(), dto.Email, dto.Name, dto.Password)

	if err != nil{
		if errors.Is(err, service.ErrEmailAlreadyExists){
			c.JSON(http.StatusConflict, gin.H{
				"error": "email já está em uso",
			})
			return
		}

		if errors.Is(err, service.ErrEmailInvalid) || 
		errors.Is(err, service.ErrPasswordTooShort) || 
		errors.Is(err, service.ErrNameTooShort){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return 
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro interno no servidor",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": new_user.ID,
		"email": new_user.Email,
		"name": new_user.Name,
	})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context){
	email_query := c.Query("email")

	if email_query == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email não pode estar vazio",
		})
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), email_query)

	if err != nil{
		if errors.Is(err, service.ErrEmailInvalid){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrUserNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"error": "usuário não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro interno no servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"nome": user.Name,
		"id": user.ID,
	})
}