package users

import (
	"practice3/internal/repository"
	"practice3/pkg/modules"
)

type Usecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *Usecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *Usecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *Usecase) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *Usecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}
