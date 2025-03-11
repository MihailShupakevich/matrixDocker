package usecase

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"matrixDocker/internal/domain"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}
func (m *MockUserRepository) FindUser(ctx context.Context, id int) (user domain.User, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id int, updateData domain.User) (updateUser domain.User, err error) {
	args := m.Called(ctx, id, updateData)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user domain.User) (newUser domain.User, err error) {
	args := m.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestFindAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := NewUserUseCase(mockRepo)
	usersData := []domain.User{
		domain.User{
			ID:       1,
			Username: "user1",
			Age:      15,
		},
		domain.User{
			ID:       2,
			Username: "user2",
			Age:      16,
		},
	}

	mockRepo.On("FindAllUsers", mock.Anything).Return(usersData, nil).Once()
	users, errFindUsers := u.FindUsers(context.Background())
	assert.NoError(t, errFindUsers)
	assert.Equal(t, usersData, users)
	mockRepo.AssertExpectations(t)
}

func TestFindUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := NewUserUseCase(mockRepo)
	userId := 1
	expectedUser := domain.User{ID: userId, Username: "user1", Age: 15}
	mockRepo.On("FindUser", mock.Anything, userId).Return(expectedUser, nil).Once()
	user, errFind := u.FindUser(context.Background(), userId)
	assert.NoError(t, errFind)
	assert.Equal(t, expectedUser.ID, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := NewUserUseCase(mockRepo)
	userData := domain.User{
		Username: "Salaga",
		Age:      14,
	}
	mockRepo.On("CreateUser", mock.Anything, userData).Return(userData, nil).Once()
	newUserCreated, errCreate := u.CreateUser(context.Background(), userData)
	assert.NoError(t, errCreate)
	assert.Equal(t, userData, newUserCreated)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := NewUserUseCase(mockRepo)
	updatedUserData := domain.User{ID: 1, Username: "Salaga", Age: 15}
	idUser := 1
	mockRepo.On("UpdateUser", mock.Anything, idUser, updatedUserData).Return(updatedUserData, nil).Once()
	_, errUpdate := u.UpdateUser(context.Background(), idUser, updatedUserData)
	assert.NoError(t, errUpdate)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := NewUserUseCase(mockRepo)
	idUser := 1
	mockRepo.On("DeleteUser", mock.Anything, idUser).Return("User  successfully deleted", nil).Once()
	success, errDelete := u.DeleteUser(context.Background(), idUser)
	assert.NoError(t, errDelete)
	assert.Equal(t, "User  successfully deleted", success)
	mockRepo.AssertExpectations(t)
}
