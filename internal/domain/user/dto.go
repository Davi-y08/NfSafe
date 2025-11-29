package user

type CreateUserDto struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginDto struct{
	Email string `json:"email"`
	PassWord string `json:"password"`
}