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
	dsn := "host=db user=dunice password=dunice dbname=dunice port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "user-topic",
		Balancer: &kafka.Hash{},
	}

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(*userRepo)
	userHandler := handlers.NewUserHandler(userUseCase, writer)

	router := gin.Default()

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.FindUsers)
		userRoutes.GET("/:id", userHandler.FindUser)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	fmt.Println("3. User routes are set up")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
