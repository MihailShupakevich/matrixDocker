package usecase

import (
	"errors"
	"matrixDocker/API/internal/domain"
	"matrixDocker/API/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

// UseCase интерфейс для работы с пользователями
type UseCase interface {
	FindUser(id int) (domain.User, error)
	FindUsers() ([]domain.User, error)
	CreateUser(newUser domain.User) (string, error)
	DeleteUser(id int) (string, error)
	UpdateUser(id int, updateUser domain.User) (domain.User, error)
}

// NewUser UseCase создает новый экземпляр UserUseCase
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

// FindUser  находит пользователя по ID
func (u *UserUseCase) FindUser(id int) (domain.User, error) {
	user, err := u.repo.FindUser(id)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (u *UserUseCase) FindUsers() ([]domain.User, error) {
	users, err := u.repo.FindAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser  создает нового пользователя
func (u *UserUseCase) CreateUser(newUser domain.User) (string, error) {
	err := u.repo.CreateUser(newUser)
	if err != nil {
		return "", err
	}
	return "User  created successfully", nil
}

// DeleteUser  удаляет пользователя по ID
func (u *UserUseCase) DeleteUser(id int) (string, error) {
	err := u.repo.DeleteUser(id)
	if err != nil {
		return "", errors.New("user not found")
	}
	return "User  deleted successfully", nil
}

// UpdateUser  обновляет данные пользователя
func (u *UserUseCase) UpdateUser(id int, updateUser domain.User) (domain.User, error) {
	existingUser, err := u.repo.FindUser(id)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	updateUser.ID = existingUser.ID // Сохраняем старый ID
	err = u.repo.UpdateUser(updateUser)
	if err != nil {
		return domain.User{}, err
	}
	return updateUser, nil
}
