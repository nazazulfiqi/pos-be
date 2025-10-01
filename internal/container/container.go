package container

import (
	"pos-be/internal/handler"
	"pos-be/internal/repository"
	"pos-be/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewContainer(db *gorm.DB) *Container {
	// repository
	userRepo := repository.NewUserRepository(db)

	// service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	// handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
