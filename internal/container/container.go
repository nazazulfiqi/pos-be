package container

import (
	"pos-be/internal/handler"
	"pos-be/internal/repository"
	"pos-be/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler          *handler.UserHandler
	AuthHandler          *handler.AuthHandler
	CategoryHandler      *handler.CategoryHandler
	ProductHandler       *handler.ProductHandler
	StockMovementHandler *handler.StockMovementHandler
}

func NewContainer(db *gorm.DB) *Container {
	// repository
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	stockMovementRepo := repository.NewStockMovementRepository(db)

	// service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	stockMovementService := service.NewStockMovementService(stockMovementRepo)

	// handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	stockMovementHandler := handler.NewStockMovementHandler(stockMovementService)

	return &Container{
		UserHandler:          userHandler,
		AuthHandler:          authHandler,
		CategoryHandler:      categoryHandler,
		ProductHandler:       productHandler,
		StockMovementHandler: stockMovementHandler,
	}
}
