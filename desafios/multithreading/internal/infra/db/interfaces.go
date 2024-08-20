package db

import "github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/entity"

type UserDAOInterface interface {
	Create(*entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
