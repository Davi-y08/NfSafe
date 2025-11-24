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
	"github.com/stripe/stripe-go/v78/customer"
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
	ErrEmailInvalid = errors.New("email inválido")
	ErrPassInvalid = errors.New("senha deve conter no mínimo 6 dígitos")
	ErrInDatabase = errors.New("ocorreu um erro com o banco de dados")
	ErrUserNotFound = errors.New("o usuário não foi encontrado")
	ErrNameInvalid = errors.New("o nome deve ter no mínimo 4 dígitos")
	ErrCriptPass = errors.New("ocorreu um erro ao criptografar a senha")
	ErrStriperError = errors.New("ocorreu um erro com o striper_id")
	ErrLogin = errors.New("email ou senha incorreta")
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
		return nil, ErrInDatabase
	}

	if u == nil{
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (s *UserService) GetUserById(ctx context.Context, id uint) (*user.User, error){
	u, err := s.repo.GetUserById(ctx, id)

	if err != nil{
		return nil, ErrInDatabase
	}

	if u == nil{
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name, password string) error {
	if len(password) < 6{
		return ErrPassInvalid
	}

	if !validateEmail(email){
		return ErrEmailInvalid
	}

	if len(name) < 4{
		return ErrNameInvalid
	}	

	existing, erro := s.repo.GetUserByEmail(ctx, email)

	if erro != nil {
		return ErrInDatabase
	}

	if existing != nil{
		return ErrEmailInvalid
	}

	hash, err := hashPassword(password)

	if err != nil{
		return ErrCriptPass
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name: stripe.String(name),
	}

	stripeCustomer, err := customer.New(params)

	if err != nil{
		return ErrStriperError
	}

	new_user := &user.User{
		Name: name,
		Email: email,
		PasswordHash: hash,
		Role: "user",
		StripeCustomerID: stripeCustomer.ID,
	}

	return s.repo.CreateUser(ctx, new_user)
}

func (s *UserService) Login(ctx context.Context, email, senha string) (*user.User, error){ 
	if email == "" || senha == ""{
		return nil, ErrLogin
	}

	u, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil{
		return nil, ErrInDatabase
	}

	if u == nil{
		return nil, ErrLogin
	}

	pass_hash := u.PasswordHash

	if !CheckPassword(pass_hash, senha){
		return nil, ErrLogin
	}

	return u, nil
}