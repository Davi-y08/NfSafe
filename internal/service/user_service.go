package service

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time" // Adicionei pra timeout
	"nf-safe/internal/domain/user"
	repo "nf-safe/internal/repository/user"
	"os"
	"github.com/stripe/stripe-go/v78"
	"golang.org/x/crypto/bcrypt"
)

func init(){
	stripe.Key = os.Getenv("STRIPE_KEY")
	if stripe.Key == "" {
		panic("STRIPE_SECRET_KEY não configurada no .env")
	}
	
	stripe.SetBackend(stripe.APIBackend, stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second, // Timeout max por chamada (não trava servidor)
		},
		MaxNetworkRetries: stripe.Int64(2), // Retry automático em falhas de rede
	}))

} 


type UserService struct{
	repo *repo.UserRepository
}

var (
	ErrEmailInvalid       = errors.New("email inválido")
	ErrEmailAlreadyExists = errors.New("email já está em uso")
	ErrPasswordTooShort   = errors.New("senha deve ter no mínimo 6 caracteres")
	ErrNameTooShort       = errors.New("nome deve ter no mínimo 4 caracteres")
	ErrDatabase           = errors.New("erro interno no banco de dados")
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrHashPassword       = errors.New("falha ao criptografar senha")
	ErrInvalidCredentials = errors.New("email ou senha incorretos")
)

func hashPassword(pass string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CheckPassword(hash, pass string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func validateEmail(str string) bool {
    _, err := mail.ParseAddress(str)
    return err == nil
}

func NewUserService(repository *repo.UserRepository) *UserService{
	return &UserService{repository}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*user.User, error){
	if !validateEmail(email){
		return nil, ErrEmailInvalid
	}

	u, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil{
		return nil, ErrDatabase
	}

	if u == nil{
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (s *UserService) GetUserById(ctx context.Context, id uint) (*user.User, error){
	u, err := s.repo.GetUserById(ctx, id)

	if err != nil{
		return nil, ErrDatabase
	}

	if u == nil{
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name, password string) (*user.User, error) {
	if len(password) < 6{
		return nil, ErrPasswordTooShort
	}

	if !validateEmail(email){
		return nil, ErrEmailInvalid
	}

	if len(name) < 4{
		return nil, ErrNameTooShort
	}	

	existing, erro := s.repo.GetUserByEmail(ctx, email)

	if erro != nil {
		return nil, ErrDatabase
	}

	if existing != nil{
		return nil, ErrEmailAlreadyExists
	}

	hash, err := hashPassword(password)

	if err != nil{
		return nil, ErrHashPassword
	}

	new_user := &user.User{
		Name: name,
		Email: email,
		PasswordHash: hash,
		Role: "user",
	}

	if err := s.repo.CreateUser(ctx, new_user); err != nil {
		return nil, ErrDatabase
	}

	return new_user, nil
}

func (s *UserService) Login(ctx context.Context, email, senha string) (*user.User, error){
	if email == "" || senha == ""{
		return nil, ErrInvalidCredentials
	}

	u, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil{
		return nil, ErrDatabase
	}

	if u == nil{
		return nil, ErrInvalidCredentials
	}

	pass_hash := u.PasswordHash

	if !CheckPassword(pass_hash, senha){
		return nil, ErrInvalidCredentials
	}

	return u, nil
}