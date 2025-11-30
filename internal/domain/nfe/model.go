package nfe

import (
	"nf-safe/internal/domain/company"
	"time"
	"gorm.io/gorm"
)

type Nfe struct {
    gorm.Model

    // Relacionamento com a empresa dona da NF
    CompanyID uint            `json:"company_id"`
    Company   company.Company `gorm:"foreignKey:CompanyID"`

    // Dados essenciais da NFe
    Chave       string    `json:"chave" gorm:"size:44;uniqueIndex"` // 44 dígitos
    Number      int       `json:"number"`                           // número da NF
    Serie       int       `json:"serie"`                            // série da NF
    Value       float64   `json:"value"`                            // valor total
    EmissaoDate time.Time `json:"emissao_date"`                     // data da emissão

    // Status SEFAZ
    Status  string `json:"status"`            // autorizada, cancelada, denegada
    Protocolo string `json:"protocolo"`       // protocolo de autorização

    // XML armazenado (opcional, mas RECOMENDADO)
    XML string `json:"xml" gorm:"type:text"`
}
