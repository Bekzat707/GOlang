package usecase

import (
    "fmt"
    "practice-7/internal/entity"
    "practice-7/internal/usecase/repo"
    "practice-7/utils"
    "github.com/google/uuid"
)

type UserUseCase struct {
    repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
    return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
    userFromDB, err := u.repo.RegisterUser(user)
    if err != nil {
        return nil, "", fmt.Errorf("register user: %w", err)
    }
    sessionID := uuid.New().String()
    return userFromDB, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, error) {
    userFromRepo, err := u.repo.LoginUser(user)
    if err != nil {
        return "", fmt.Errorf("User From Repo: %w", err)
    }
    if !utils.CheckPassword(userFromRepo.Password, user.Password) {
        return "", fmt.Errorf("Check Password: %w", err)
    }
    token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
    if err != nil {
        return "", fmt.Errorf("Generate JWT: %w", err)
    }
    return token, nil
}

func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
    return u.repo.GetUserByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
    return u.repo.PromoteUser(id)
}
