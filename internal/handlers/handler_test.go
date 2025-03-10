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
	"strconv"
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

func TestFindUsers(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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

	jsonDataUser, err := json.Marshal(usersData)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(jsonDataUser))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUseCase.On("FindUsers", mock.Anything).Return(usersData, nil).Once()
	h.FindUsers(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(jsonDataUser), w.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestFindUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	userId := 1
	expectedUser := domain.User{ID: userId, Username: "user1", Age: 15}
	mockUseCase.On("FindUser", mock.Anything, userId).Return(expectedUser, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", userId), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(userId)}}

	h.FindUser(ctx)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"ID":1,"Username":"user1","Age":15}`, rr.Body.String())
	mockUseCase.AssertExpectations(t)
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
	jsonData, err := json.Marshal(userData)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	mockUseCase.On("CreateUser", mock.Anything, userData).Return(userData, nil).Once()
	h.CreateUser(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
	expectedResponse := gin.H{"newUser": userData}
	expectedJSON, err1 := json.Marshal(expectedResponse)
	assert.NoError(t, err1)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	updatedUserData := domain.User{ID: 1, Username: "Salaga", Age: 15}
	jsonData, err := json.Marshal(updatedUserData)
	assert.NoError(t, err)
	idUser := 1
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/%d", idUser), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(idUser)}}
	mockUseCase.On("UpdateUser", mock.Anything, idUser, updatedUserData).Return(updatedUserData, nil).Once()
	h.UpdateUser(ctx)
	updatedUser := new(domain.User)
	_ = ctx.ShouldBindJSON(&updatedUser)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(jsonData), w.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	idUser := 1
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%d", idUser), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(idUser)}}
	mockUseCase.On("DeleteUser", mock.Anything, idUser).Return("User deleted", nil).Once()
	h.DeleteUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"User deleted"}`, w.Body.String())
	mockUseCase.AssertExpectations(t)
}
