package service

import (
	"errors"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
 
var email_admin string
var jwt_key []byte

func init(){
	jwt_key_aux := os.Getenv("JWT_SECRET_KEY")
	email_admin_aux := os.Getenv("ADMIN_EMAIL")

	if jwt_key_aux == "" || email_admin_aux == ""{
		panic("não foi possivel acessar o .env")
	}

	jwt_key = []byte(jwt_key_aux)
	email_admin = email_admin_aux
}

func getRoleUser(email string) string{
	if email == email_admin{
		return "admin"
	}

	return "user"
}

func GenerateTokenJwt(email, username string, id uint) (*string, error){
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"preferred_username": username,
		"role": getRoleUser(email),
		"email": email,
		"iss": "minha-api",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(), 
	})

	tokenString, err := claims.SignedString(jwt_key)

	if err != nil{
		return nil, errors.New("não foi possivel gerar o token")
	}

	return &tokenString, nil
}