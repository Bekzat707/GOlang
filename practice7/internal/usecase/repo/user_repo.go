package repo

import (
    "fmt"
    "practice-7/internal/entity"
    "practice-7/pkg/postgres"
)

type UserRepo struct {
    PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
    return &UserRepo{pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
    err := u.PG.Conn.Create(user).Error
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
    var userFromDB entity.User
    if err := u.PG.Conn.Where("username = ?", user.Username).First(&userFromDB).Error; err != nil {
        return nil, fmt.Errorf("Username Not Found: %v", err)
    }
    return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(id string) (*entity.User, error) {
    var userFromDB entity.User
    if err := u.PG.Conn.Where("id = ?", id).First(&userFromDB).Error; err != nil {
        return nil, fmt.Errorf("User Not Found: %v", err)
    }
    return &userFromDB, nil
}

func (u *UserRepo) PromoteUser(id string) error {
    return u.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", "admin").Error
}
