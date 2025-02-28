package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"matrixDocker/internal/domain"
	"matrixDocker/internal/handlers/mocks"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserHandler := mocks.NewMockUserHandlerInterface(ctrl)
	ctx, _ := gin.CreateTestContext(nil)
	userData := map[string]interface{}{"username": "Salaga", "age": 14}
	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUserHandler.EXPECT().CreateUser(gomock.Any()).Times(1).Do(func(ctx *gin.Context) {
		var receivedUser domain.User
		err := ctx.ShouldBindJSON(&receivedUser)
		assert.NoError(t, err)
		assert.Equal(t, "Salaga", receivedUser.Username)
		assert.Equal(t, 14, receivedUser.Age)
	})
	mockUserHandler.CreateUser(ctx)
	ctrl.Finish()
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserHandler := mocks.NewMockUserHandlerInterface(ctrl)
	ctx, _ := gin.CreateTestContext(nil)

	var userData domain.User
	userData.Username = "Salaga"
	userData.Age = 14

	userUpdateData := map[string]interface{}{
		"username": "Salaga",
		"age":      15,
	}
	jsonData, _ := json.Marshal(userUpdateData)
	req, _ := http.NewRequest("PATCH", "/id", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	mockUserHandler.EXPECT().UpdateUser(gomock.Any()).Times(1).Do(func(ctx *gin.Context) {
		var updateUser domain.User
		err := ctx.ShouldBindJSON(&updateUser)
		assert.Equal(t, userData.Username, updateUser.Username)
		assert.NotEqual(t, userData.Age, updateUser.Age)
		assert.NoError(t, err)
		assert.Equal(t, "Salaga", updateUser.Username)
		assert.Equal(t, 15, updateUser.Age)
	})

	mockUserHandler.UpdateUser(ctx)
	ctrl.Finish()
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserHandler := mocks.NewMockUserHandlerInterface(ctrl)
	ctx, _ := gin.CreateTestContext(nil)
	idUser := 1
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%d", idUser), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUserHandler.EXPECT().DeleteUser(gomock.Any()).Times(1).Do(func(ctx *gin.Context) {})
	mockUserHandler.DeleteUser(ctx)
	ctrl.Finish()
}

func TestFindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserHandler := mocks.NewMockUserHandlerInterface(ctrl)
	ctx, _ := gin.CreateTestContext(nil)
	idUser := 1
	req, _ := http.NewRequest("GET", fmt.Sprintf("%d", idUser), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	mockUserHandler.EXPECT().FindUser(gomock.Any()).Times(1).Do(func(ctx *gin.Context) {})
	mockUserHandler.FindUser(ctx)
	ctrl.Finish()
}
func TestFindUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserHandler := mocks.NewMockUserHandlerInterface(ctrl)
	ctx, _ := gin.CreateTestContext(nil)

	usersData := []domain.User{
		domain.User{1, "alex", 14},
		domain.User{2, "bob", 15},
		domain.User{3, "david", 17},
	}
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	mockUserHandler.EXPECT().FindUsers(gomock.Any()).Times(1).Do(func(ctx *gin.Context) {

		receivedUsers := []domain.User{
			domain.User{1, "alex", 14},
			domain.User{2, "bob", 15},
			domain.User{3, "david", 17},
		}
		assert.Equal(t, usersData, receivedUsers)

	})
	mockUserHandler.FindUsers(ctx)
	ctrl.Finish()
}
