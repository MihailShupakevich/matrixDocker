package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"matrixDocker/API/internal/domain"
	"matrixDocker/API/internal/handlers"
	"matrixDocker/API/internal/repository"
	"matrixDocker/API/internal/usecase"
)

func main() {

	dsn := "host=localhost user=postgres password=admin dbname=matrixDocker port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Инициализация репозитория, usecase и обработчика
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(*userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Настройка маршрутов
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

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
