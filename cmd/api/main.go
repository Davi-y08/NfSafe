package main

import (
	"log"
	"net/http"
	"nf-safe/internal/domain/user"
	"nf-safe/internal/infra/db"
	repository "nf-safe/internal/repository/user"

	"github.com/gin-gonic/gin"
)

func main() {

	database := db.Connect()

	// 1. Rodar migrações ANTES de tudo
	db.RunMigrations(database)

	// 2. Criar repo
	repo := repository.NewUserRepository(database)

	// 3. Testar inserção
	u := &user.User{
		Name:         "Elton Davi",
		Email:        "teste@gmail.com",
		PasswordHash: "elton123",
	}

	err := repo.CreateUser(u)
	if err != nil {
		log.Println("erro ao criar:", err)
	} else {
		log.Println("usuário criado! ID:", u.ID)
	}

	// 4. Subir servidor
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("API rodando em http://localhost:8080")
	r.Run(":8080")
}
