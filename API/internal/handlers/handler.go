package handlers

import (
	"github.com/gin-gonic/gin"
	"matrixDocker/API/internal/domain"
	"matrixDocker/API/internal/usecase"
	"net/http"
	"strconv"
)

type UserHandler struct {
	useCase *usecase.UserUseCase
}

type UserHandlerInterface interface {
	FindUser(c *gin.Context)
	FindUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

func NewUserHandler(usecase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: usecase,
	}
}

func (h *UserHandler) FindUsers(ctx *gin.Context) {
	allUsers, err := h.useCase.FindUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding users"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": allUsers})
}

func (h *UserHandler) FindUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User  ID is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.useCase.FindUser(idInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User  not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	updateUser := new(domain.User)
	err := ctx.BindJSON(updateUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User  ID is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.useCase.UpdateUser(idInt, *updateUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User  ID is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	message, err := h.useCase.DeleteUser(idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	body := new(domain.User)
	err := ctx.BindJSON(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.useCase.CreateUser(*body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"newUser": newUser})
}
