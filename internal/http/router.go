package http

import (
	"nf-safe/internal/http/handlers"
	"nf-safe/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(userService *service.UserService) *gin.Engine{
	r := gin.Default()

	userHandler := handlers.NewUserHandler(userService)
	
}