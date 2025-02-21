package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"matrixDocker/internal/domain"
	"matrixDocker/internal/handlers"
	"matrixDocker/internal/repository"
	"matrixDocker/internal/usecase"
)

func main() {
	dsn := "host=localhost user=dunice password=dunice dbname=dunice port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("1. Database is active!")

	// Автоматическая миграция структуры User
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("2. User entity migrated to the database")

	// Создание Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), // Укажите адрес вашего Kafka брокера
		Topic:    "user-topic",                // Укажите тему, в которую будете отправлять сообщения
		Balancer: &kafka.Hash{},
	}

	// Инициализация репозитория, юзкейса и обработчика
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(*userRepo) // Передаем указатель на userRepo
	userHandler := handlers.NewUserHandler(userUseCase, writer)

	// Настройка маршрутизатора Gin
	router := gin.Default()

	// Определение маршрутов для пользователей
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.FindUsers)
		userRoutes.GET("/:id", userHandler.FindUser)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	fmt.Println("3. User routes are set up")

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
