package handlers

import (
	"errors"
	"net/http"
	"nf-safe/internal/domain/user"
	"nf-safe/internal/service"
	"strconv"

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
		return
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

func (h *UserHandler) GetUserById(c *gin.Context){
	id_string := c.Query("id")

	if id_string == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "o id não pode estar vazio",
		})
	}

	id, err := strconv.ParseUint(id_string, 10, 32)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "erro ao converter a id para Uint",
		})
		return 
	}

	var id_int uint = uint(id)

	user, err := h.userService.GetUserById(c.Request.Context(), id_int)

	if err != nil{
		if errors.Is(err, service.ErrEmailInvalid){
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

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"nome": user.Name,
		"id": user.ID,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var dto user.LoginDto

	if err := c.ShouldBindJSON(&dto); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": "corpo inválido",
		})
		return 
	}

	user, err := h.userService.Login(c.Request.Context(), dto.Email, dto.PassWord)

	if err != nil{
		if errors.Is(err, service.ErrDatabase){
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrInvalidCredentials){
			c.JSON(http.StatusBadRequest, gin.H{
				"error:": "email ou senha incorreta",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro interno no servidor",
		})
		return
	}

	if user == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "usuário não encontrado",
		})
		return
	}

	token, err := service.GenerateTokenJwt(user.Email, user.Name, user.ID)

	if err != nil || token == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "erro ao gerar token jwt",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

