package db

import (
    "gorm.io/gorm"
    "nf-safe/internal/domain/user"
    "nf-safe/internal/domain/company"
    "nf-safe/internal/domain/nfe"
)

func RunMigrations(db *gorm.DB) {
    db.AutoMigrate(&user.User{})
    db.AutoMigrate(&company.Company{})
    db.AutoMigrate(&nfe.Nfe{})
}
