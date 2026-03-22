package database

import "github.com/renamrgb/go-expert-apis/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}


