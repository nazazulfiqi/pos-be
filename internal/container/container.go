package container

import (
	"pos-be/internal/handler"
	"pos-be/internal/repository"
	"pos-be/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler     *handler.UserHandler
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
}

func NewContainer(db *gorm.DB) *Container {
	// repository
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	return &Container{
		UserHandler:     userHandler,
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
	}
}
