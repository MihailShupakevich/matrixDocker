package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"io/ioutil"
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

	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUseCase.On("FindUsers", ctx).Return(usersData, nil).Once()
	h.FindUsers(ctx)
	assert.Equal(t, nil, req.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	mockUseCase.AssertExpectations(t)
}

func TestFindUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	userId := 1
	expectedUser := domain.User{ID: userId, Username: "user1", Age: 15}
	mockUseCase.On("FindUser", ctx, userId).Return(expectedUser, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", userId), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(userId)}}
	h.FindUser(ctx)
	assert.Equal(t, nil, req.Body)
	assert.NotNil(t, req)
	assert.NotNil(t, ctx.Params)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, http.MethodGet, req.Method)
	idParam, errConv := strconv.Atoi(ctx.Param("id"))
	assert.NoError(t, errConv)
	assert.Equal(t, idParam, userId)
	assert.JSONEq(t, `{"ID":1,"Username":"user1","Age":15}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	mockUseCase.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "user-topic",
		Balancer: &kafka.Hash{},
	}
	h := NewUserHandler(mockUseCase, kafkaWriter)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	userData := domain.User{
		Username: "Salaga",
		Age:      14,
	}
	jsonData, err := json.Marshal(userData)
	assert.Nil(t, err)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUseCase.On("CreateUser", ctx, userData).Return(userData, nil).Once()
	h.CreateUser(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, req.Body, nil)
	assert.NotEqual(t, ctx, nil)
	expectedResponse := gin.H{"newUser": userData}
	expectedJSON, err1 := json.Marshal(expectedResponse)
	assert.NoError(t, err1)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())
	assert.Equal(t, http.MethodPost, req.Method)
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
	mockUseCase.On("UpdateUser", ctx, idUser, updatedUserData).Return(updatedUserData, nil).Once()
	h.UpdateUser(ctx)
	updatedUser := new(domain.User)
	_ = ctx.ShouldBindJSON(&updatedUser)
	assert.NotEqual(t, req.Body, nil)
	idParam, errConv := strconv.Atoi(ctx.Param("id"))
	assert.NoError(t, errConv)
	assert.Equal(t, idParam, idUser)
	assert.Equal(t, http.MethodPatch, req.Method)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
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
	mockUseCase.On("DeleteUser", ctx, idUser).Return("User deleted", nil).Once()
	h.DeleteUser(ctx)
	assert.Equal(t, req.Body, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, req.Method, http.MethodDelete)
	idParam, errConv := strconv.Atoi(ctx.Param("id"))
	assert.NoError(t, errConv)
	assert.Equal(t, idParam, idUser)
	assert.JSONEq(t, `{"message":"User deleted"}`, w.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestFindUser_Errors(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		statusCode   int
		errorMessage string
	}{
		{
			name:         "Missing ID",
			paramID:      "",
			statusCode:   http.StatusBadRequest,
			errorMessage: "User ID is required",
		},
		{
			name:         "Invalid ID",
			paramID:      "abc",
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid user ID",
		},
		{name: "User Not Found",
			paramID:      "1",
			statusCode:   http.StatusNotFound,
			errorMessage: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			h := NewUserHandler(mockUseCase, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			if tt.statusCode == http.StatusNotFound {
				mockUseCase.On("FindUser", mock.Anything, 1).Return(domain.User{}, errors.New("user not found")).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.paramID), nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			ctx.Params = gin.Params{gin.Param{Key: "id", Value: tt.paramID}}
			h.FindUser(ctx)
			assert.Equal(t, tt.statusCode, w.Code)
			responseBody, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(responseBody), tt.errorMessage)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestFindUsers_Errors(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		errorMessage string
	}{
		{
			name:         "Error finding",
			statusCode:   http.StatusInternalServerError,
			errorMessage: "Error finding users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			h := NewUserHandler(mockUseCase, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			if tt.statusCode == http.StatusInternalServerError {
				mockUseCase.On("FindUsers", mock.Anything).Return([]domain.User{}, errors.New("user not found")).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			h.FindUsers(ctx)
			assert.Equal(t, tt.statusCode, w.Code)
			responseBody, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(responseBody), tt.errorMessage)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestCreateUser_Errors(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		errorMessage string
		jsonData     []byte
	}{
		{
			name:         "BindJSON",
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid JSON body",
			jsonData:     []byte("invalid json"),
		},
		{
			name:         "User not created",
			statusCode:   http.StatusInternalServerError,
			errorMessage: "Error creating user",
			jsonData:     []byte(`{"username": "Salaga", "age": 14}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			h := NewUserHandler(mockUseCase, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			if tt.name == "User not created" {
				mockUseCase.On("CreateUser", mock.Anything, mock.Anything).Return(domain.User{}, errors.New("user not created")).Once()
			}

			h.CreateUser(ctx)
			assert.Equal(t, tt.statusCode, w.Code)
			responseBody, errRead := ioutil.ReadAll(w.Body)
			assert.NoError(t, errRead)
			assert.Contains(t, string(responseBody), tt.errorMessage)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUpdateUser_Errors(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		statusCode   int
		errorMessage string
		jsonData     []byte
	}{
		{
			name:         "BindJSON",
			paramID:      "1",
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid JSON body",
			jsonData:     []byte("invalid json"),
		},
		{
			name:         "ctx_Param",
			paramID:      "",
			statusCode:   http.StatusBadRequest,
			errorMessage: "User ID is required",
			jsonData:     []byte(`{"username": "Salaga", "age": 15}`),
		},
		{
			name:         "conv_Atoi",
			paramID:      "abc",
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid user ID",
			jsonData:     []byte(`{"username": "Salaga", "age": 15}`),
		},
		{
			name:         "Can't update",
			paramID:      "0",
			statusCode:   http.StatusInternalServerError,
			errorMessage: "Error updating user",
			jsonData:     []byte(`{"Username": "Salaga", "Age": 15}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			h := NewUserHandler(mockUseCase, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{gin.Param{Key: "id", Value: tt.paramID}}
			req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s", tt.paramID), bytes.NewBuffer(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			if tt.paramID == "" {
				h.UpdateUser(ctx)
				assert.Equal(t, tt.statusCode, w.Code)
				responseBody, errRead := ioutil.ReadAll(w.Body)
				assert.NoError(t, errRead)
				assert.Contains(t, string(responseBody), tt.errorMessage)
				return
			}
			if tt.paramID == "abc" {
				h.UpdateUser(ctx)
				assert.Equal(t, tt.statusCode, w.Code)
				responseBody, errRead := ioutil.ReadAll(w.Body)
				assert.NoError(t, errRead)
				assert.Contains(t, string(responseBody), tt.errorMessage)
				return
			}
			id, err := strconv.Atoi(tt.paramID)
			if err != nil {
				t.Fatalf("Failed to convert paramID to int: %v", err)
			}

			if tt.name == "Can't update" {
				mockUseCase.On("UpdateUser", mock.Anything, id, mock.Anything).Return(domain.User{}, errors.New("user not updated")).Once()
			}

			h.UpdateUser(ctx)
			assert.Equal(t, tt.statusCode, w.Code)
			responseBody, errRead := ioutil.ReadAll(w.Body)
			assert.NoError(t, errRead)
			assert.Contains(t, string(responseBody), tt.errorMessage)
			mockUseCase.AssertExpectations(t)

		})
	}
}

func TestDeleteUser_Errors(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		statusCode   int
		errorMessage string
	}{
		{
			name:         "ctx_Param",
			paramID:      "",
			statusCode:   http.StatusBadRequest,
			errorMessage: "User  ID is required",
		},
		{
			name:         "conv_Atoi",
			paramID:      "abc",
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid user ID",
		},
		{
			name:         "Can't delete",
			paramID:      "1",
			statusCode:   http.StatusInternalServerError,
			errorMessage: "Error deleting user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			h := NewUserHandler(mockUseCase, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{gin.Param{Key: "id", Value: tt.paramID}}
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s", tt.paramID), nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			if tt.paramID == "" {
				h.DeleteUser(ctx)
				assert.Equal(t, tt.statusCode, w.Code)
				responseBody, errRead := ioutil.ReadAll(w.Body)
				assert.NoError(t, errRead)
				assert.Contains(t, string(responseBody), tt.errorMessage)
				return
			}
			if tt.paramID == "abc" {
				h.DeleteUser(ctx)
				assert.Equal(t, tt.statusCode, w.Code)
				responseBody, errRead := ioutil.ReadAll(w.Body)
				assert.NoError(t, errRead)
				assert.Contains(t, string(responseBody), tt.errorMessage)
				return
			}
			id, err := strconv.Atoi(tt.paramID)
			if err != nil {
				t.Fatalf("Failed to convert paramID to int: %v", err)
			}
			if tt.name == "Can't delete" {
				mockUseCase.On("DeleteUser", mock.Anything, id).Return("", errors.New("user not deleted")).Once()
			}
			h.DeleteUser(ctx)
			assert.Equal(t, tt.statusCode, w.Code)
			responseBody, errRead := ioutil.ReadAll(w.Body)
			assert.NoError(t, errRead)
			assert.Contains(t, string(responseBody), tt.errorMessage)
			mockUseCase.AssertExpectations(t)
		})
	}
}
