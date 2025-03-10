package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"log"
	"matrixDocker/internal/domain"
	"matrixDocker/internal/usecase"
	"net/http"
	"strconv"
)

type UserHandlerInterface interface {
	FindUsers(ctx *gin.Context)
	FindUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
}

type KafkaMessage struct {
	ID   int         `json:"id"`
	Data domain.User `json:"data"`
}

type UserHandler struct {
	useCase usecase.UseCaseI
	writer  *kafka.Writer
}

func NewUserHandler(usecase usecase.UseCaseI, writer *kafka.Writer) *UserHandler {
	return &UserHandler{
		useCase: usecase,
		writer:  writer,
	}
}

func (h *UserHandler) FindUsers(ctx *gin.Context) {
	allUsers, err := h.useCase.FindUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding users"})
		return
	}
	ctx.JSON(http.StatusOK, allUsers)
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

	user, err := h.useCase.FindUser(ctx, idInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User  not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
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

	user, err := h.useCase.UpdateUser(ctx, idInt, *updateUser)
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

	message, err := h.useCase.DeleteUser(ctx, idInt)
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

	newUser, err := h.useCase.CreateUser(ctx, *body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Отправка сообщения в Kafka
	go h.sendToKafka(newUser)

	ctx.JSON(http.StatusCreated, gin.H{"newUser ": newUser})
}

func (h *UserHandler) sendToKafka(user domain.User) {
	kafkaMessage := KafkaMessage{
		ID:   user.ID, // Предполагается, что у вас есть поле ID в структуре User
		Data: user,
	}

	msg, err := json.Marshal(kafkaMessage)
	if err != nil {
		log.Println("Error marshaling user:", err)
		return
	}

	err = h.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprint(kafkaMessage.ID)),
			Value: msg,
		},
	)
	if err != nil {
		log.Println("Error sending message to Kafka:", err)
	}
}
