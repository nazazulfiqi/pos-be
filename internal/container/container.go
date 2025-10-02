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
	ProductHandler  *handler.ProductHandler
}

func NewContainer(db *gorm.DB) *Container {
	// repository
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)

	// handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)

	return &Container{
		UserHandler:     userHandler,
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
	}
}
