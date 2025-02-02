package repository

import (
	"gorm.io/gorm"
	"matrixDocker/API/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindUser(userId int) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, userId).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(userId int, updateUser domain.User) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userId).Error; err != nil {
		return domain.User{}, err
	}
	if err := r.db.Model(&user).Updates(updateUser).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id int) (string, error) {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return "", err
	}
	return "User  successfully deleted", nil
}

func (r *UserRepository) CreateUser(newUser domain.User) (string, error) {
	if err := r.db.Create(&newUser).Error; err != nil {
		return "", err
	}
	return "User  successfully created", nil
}
