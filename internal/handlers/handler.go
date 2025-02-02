package handlers

import "matrixDocker/internal/usecase"

type UserHandler struct {
	usecase *usecase.Usecase
}

type UserHandlerInterface interface {
	findUser(c *ginContext)
}
