package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex;size:180"`
	PasswordHash string `json:"-"`

	Plan string `json:"plan"`
	PlanExpiresAt *time.Time `json:"plan_expires_at"`
	TrialEndsAt *time.Time `json:"trials_end_at"`
	IsActive bool `json:"is_active" gorm:"default:true"`

	StripeCustomerID string `json:"stripe_customer_id" gorm:"uniqueIndex;size:60"`
	Role string `json:"role" gorm:"default:user"`

	SubscriptionStatus string `json:"subscription_status" gorm:"default:'none'"`
}

