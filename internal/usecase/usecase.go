package usecase

import (
	"errors"
	"matrixDocker/internal/domain"
	"matrixDocker/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

type UseCase interface {
	FindUser(id int) (domain.User, error)
	FindUsers() ([]domain.User, error)
	CreateUser(newUser domain.User) (domain.User, error)
	DeleteUser(id int) (string, error)
	UpdateUser(id int, updateUser domain.User) (domain.User, error)
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

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

func (u *UserUseCase) CreateUser(newUser domain.User) (domain.User, error) {
	newUser, err := u.repo.CreateUser(newUser)
	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (u *UserUseCase) DeleteUser(id int) (string, error) {
	_, err := u.repo.DeleteUser(id)
	if err != nil {
		return "", errors.New("user not found")
	}
	return "User  deleted successfully", nil
}

func (u *UserUseCase) UpdateUser(id int, updateUser domain.User) (domain.User, error) {
	_, err := u.repo.FindUser(id)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	_, err = u.repo.UpdateUser(id, updateUser)
	if err != nil {
		return domain.User{}, err
	}
	return updateUser, nil
}
