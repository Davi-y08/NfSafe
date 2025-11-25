package handlers

import (
	"net/http"
	"nf-safe/internal/domain/user"
	"nf-safe/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

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

	ctx := c.Request.Context()
	
	err := h.userService.CreateUser(ctx, dto.Email, dto.Name, dto.Password)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, err)
}