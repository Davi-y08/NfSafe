package company

type CreateCompanyDto struct{
	UserID string `json:"user_id"`
	Cnpj string `json:"cnpj"`
	Name string `json:"name"`
	RazaoSocial string `json:"razao_social"`
	NomeFantasia string `json:"nome_fantasia"`
	Status string `json:"status"`
}

