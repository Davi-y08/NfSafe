package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "log"
)

func Connect() *gorm.DB{
	dns := "host=localhost user=postgres password=eltinho123 dbname=nfsafe port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil{
		log.Fatalf("erro ao se conectar com o banco de dados: %v", err)
	}

	return db
}