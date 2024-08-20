package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UID = uuid.UUID

func newId() UID {
	return UID(uuid.New())
}

type User struct {
	ID       UID    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) ParseID(s string) (UID, error) {
	id, err := uuid.Parse(s)
	return UID(id), err
}

func NewUser(email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       newId(),
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) == nil
}
