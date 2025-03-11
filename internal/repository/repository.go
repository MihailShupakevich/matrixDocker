package repository

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"matrixDocker/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryI interface {
	FindAllUsers(ctx context.Context) ([]domain.User, error)
	FindUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, id int, updateUser domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) (string, error)
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindUser(ctx context.Context, userId int) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, userId).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userId int, updateUser domain.User) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userId).Error; err != nil {
		return domain.User{}, err
	}
	if err := r.db.Model(&user).Updates(updateUser).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) (string, error) {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return "", err
	}
	return "User successfully deleted", nil
}

func (r *UserRepository) CreateUser(ctx context.Context, newUser domain.User) (domain.User, error) {
	if err := r.db.Create(&newUser).Error; err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}
