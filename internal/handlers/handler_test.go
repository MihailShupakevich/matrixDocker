package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"matrixDocker/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) FindUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}
func (m *MockUserUseCase) FindUser(ctx context.Context, id int) (user domain.User, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(ctx context.Context, id int, updateData domain.User) (updateUser domain.User, err error) {
	args := m.Called(ctx, id, updateData)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id int) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockUserUseCase) CreateUser(ctx context.Context, user domain.User) (newUser domain.User, err error) {
	args := m.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {

	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	userData := domain.User{
		Username: "Salaga",
		Age:      14,
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	mockUseCase.On("CreateUser", mock.Anything, userData).Return(userData, nil).Once()
	h.CreateUser(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	userData := domain.User{
		ID:       0,
		Username: "Salaga",
		Age:      14,
	}
	updatedUserData := domain.User{ID: 0, Username: "Salaga", Age: 15}
	jsonData, err := json.Marshal(updatedUserData)
	fmt.Println(err)
	idUser := 0
	req, _ := http.NewRequest("PATCH", "/0", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUseCase.On("UpdateUser", mock.Anything, idUser, updatedUserData).Return(idUser, updatedUserData, nil).Once()
	h.UpdateUser(ctx)
	assert.Equal(t, userData.Username, updatedUserData.Username)
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	idUser := 1
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%d", idUser), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUseCase.On("DeleteUser ", mock.Anything, idUser).Return("User  deleted", nil).Once()
	h.DeleteUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestFindUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	idUser := 1
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", idUser), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	mockUseCase.On("FindUser", mock.Anything).Return(domain.User{
		Username: "Alex",
		Age:      20,
	}, nil).Once()
	h.FindUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	response := new(domain.User)
	err := ctx.ShouldBindJSON(response)
	assert.Nil(t, err)
	assert.Equal(t, "Alex", response.Username)
	assert.Equal(t, 20, response.Age)
	mockUseCase.AssertExpectations(t)
}

func TestFindUsers(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	mockUseCase.On("FindUsers", mock.Anything).Return([]domain.User{
		{
			Username: "Alex",
			Age:      20,
		},
		{
			Username: "Bob",
			Age:      25,
		},
	}, nil).Once()
	h.FindUsers(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	var response domain.User
	err := ctx.BindJSON(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Alex", response.Username)
	assert.Equal(t, 20, response.Age)
	assert.Equal(t, "Bob", response.Username)
	assert.Equal(t, 25, response.Age)
	mockUseCase.AssertExpectations(t)
}
