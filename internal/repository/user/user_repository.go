package repository

import (
	"context"
	"errors"
	"nf-safe/internal/domain/user"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error{
	if u == nil{
		return errors.New("não foi possivel criar usuário")
	}

	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error){
	if len(email) == 0{
		return nil, errors.New("o email está vazio")
	}
	
	var u user.User

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &u, nil
}

func (r *UserRepository) UpdateUserPlan(ctx context.Context, userId uint, plan string, expiresAt *time.Time) error {
    return r.db.Model(&user.User{}).WithContext(ctx).
        Where("id = ?", userId).
        Updates(&user.User{
			Plan: plan,
			PlanExpiresAt: expiresAt,
			SubscriptionStatus: "active",
		}).Error
}

func (r *UserRepository) UpdateSubscriptionUser(ctx context.Context, userId uint, status string) error{
	return r.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", userId).UpdateColumn("subscription_status", status).Error
}

func (r *UserRepository) GetUserById(ctx context.Context, id uint) (*user.User, error){
	var u user.User

	result := r.db.WithContext(ctx).First(&u, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &u, nil
}