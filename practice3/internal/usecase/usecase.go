package usecase

import (
	"practice3/internal/repository"
	"practice3/internal/usecase/users"
	"practice3/pkg/modules"
)

type UserUsecase interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) error
}

type Usecases struct {
	UserUsecase
}

func NewUsecases(repos *repository.Repositories) *Usecases {
	return &Usecases{
		UserUsecase: users.NewUserUsecase(repos.UserRepository),
	}
}
