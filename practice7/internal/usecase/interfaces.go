package usecase

import (
    "practice-7/internal/entity"
)

type UserInterface interface {
    RegisterUser(user *entity.User) (*entity.User, string, error)
    LoginUser(user *entity.LoginUserDTO) (string, error)
    GetUserByID(id string) (*entity.User, error)
    PromoteUser(id string) error
}
