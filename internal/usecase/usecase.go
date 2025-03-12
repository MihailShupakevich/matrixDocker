package usecase

import (
	"errors"
	"golang.org/x/net/context"
	"matrixDocker/internal/domain"
	"matrixDocker/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepositoryI
}

type UseCaseI interface {
	FindUsers(ctx context.Context) ([]domain.User, error)
	FindUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, id int, updateUser domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) (string, error)
	CreateUser(ctx context.Context, newUser domain.User) (domain.User, error)
}

func NewUserUseCase(repo repository.UserRepositoryI) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) FindUser(ctx context.Context, id int) (domain.User, error) {
	user, err := u.repo.FindUser(ctx, id)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (u *UserUseCase) FindUsers(ctx context.Context) ([]domain.User, error) {
	users, err := u.repo.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserUseCase) CreateUser(ctx context.Context, newUser domain.User) (domain.User, error) {
	newUserCreated, err := u.repo.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, errors.New("user not created")
	}
	return newUserCreated, nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id int) (string, error) {
	_, err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return "", errors.New("user not deleted")
	}
	return "User  successfully deleted", nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, id int, updateUser domain.User) (domain.User, error) {
	_, err := u.repo.UpdateUser(ctx, id, updateUser)
	if err != nil {
		return domain.User{}, errors.New("user not updated")
	}
	return updateUser, nil
}
