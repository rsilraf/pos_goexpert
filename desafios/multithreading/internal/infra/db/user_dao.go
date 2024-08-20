package db

import (
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/entity"
	"gorm.io/gorm"
)

type UserDAO struct {
	DB *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db}
}
func (u *UserDAO) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDAO) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}
