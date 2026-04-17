package service

import (
	"errors"
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}

	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}

	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)

	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	t.Run("User already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		email := "test@test.com"
		existingUser := &repository.User{ID: 1, Name: "Existing User"}
		userToRegister := &repository.User{ID: 2, Name: "New User"}

		mockRepo.EXPECT().GetByEmail(email).Return(existingUser, nil)

		err := userService.RegisterUser(userToRegister, email)

		assert.Error(t, err)
		assert.Equal(t, "user with this email already exists", err.Error())
	})

	t.Run("New User -> Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		email := "test@test.com"
		userToRegister := &repository.User{ID: 2, Name: "New User"}

		mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(userToRegister).Return(nil)

		err := userService.RegisterUser(userToRegister, email)

		assert.NoError(t, err)
	})

	t.Run("Repository error on get by email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		email := "test@test.com"
		userToRegister := &repository.User{ID: 2, Name: "New User"}

		mockRepo.EXPECT().GetByEmail(email).Return(nil, errors.New("db error"))

		err := userService.RegisterUser(userToRegister, email)

		assert.Error(t, err)
		assert.Equal(t, "error getting user with this email", err.Error())
	})

	t.Run("Repository error on CreateUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		email := "test@test.com"
		userToRegister := &repository.User{ID: 2, Name: "New User"}

		mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(userToRegister).Return(errors.New("db error"))

		err := userService.RegisterUser(userToRegister, email)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUpdateUserName(t *testing.T) {
	t.Run("Empty name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		err := userService.UpdateUserName(1, "")

		assert.Error(t, err)
		assert.Equal(t, "name cannot be empty", err.Error())
	})

	t.Run("User not found / repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		mockRepo.EXPECT().GetUserByID(1).Return(nil, errors.New("not found"))

		err := userService.UpdateUserName(1, "New Name")

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("Successful update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		user := &repository.User{ID: 1, Name: "Old Name"}
		
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		
		mockRepo.EXPECT().UpdateUser(user).DoAndReturn(func(u *repository.User) error {
			assert.Equal(t, "New Name", u.Name)
			return nil
		})

		err := userService.UpdateUserName(1, "New Name")

		assert.NoError(t, err)
	})

	t.Run("UpdateUser Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		user := &repository.User{ID: 1, Name: "Old Name"}

		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(errors.New("db update error"))

		err := userService.UpdateUserName(1, "New Name")

		assert.Error(t, err)
		assert.Equal(t, "db update error", err.Error())
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Attempt to delete admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		err := userService.DeleteUser(1)

		assert.Error(t, err)
		assert.Equal(t, "it is not allowed to delete admin user", err.Error())
	})

	t.Run("Successful delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		mockRepo.EXPECT().DeleteUser(2).Return(nil)

		err := userService.DeleteUser(2)

		assert.NoError(t, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)

		mockRepo.EXPECT().DeleteUser(2).Return(errors.New("db error"))

		err := userService.DeleteUser(2)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
