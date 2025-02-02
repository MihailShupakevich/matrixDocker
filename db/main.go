package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

	userRepo := repository.NewUser Repository(db)

	router := gin.Default()

	// Определение маршрутов для пользователей
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userRepo.FindAllUsers) // Реализуйте методы в репозитории
		userRoutes.GET("/:id", userRepo.FindUser )
		userRoutes.POST("/", userRepo.CreateUser )
		userRoutes.PUT("/:id", userRepo.UpdateUser )
		userRoutes.DELETE("/:id", userRepo.DeleteUser )
	}

	// Запуск сервера
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
