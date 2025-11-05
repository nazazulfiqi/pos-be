package router

import (
	"pos-be/internal/container"
	"pos-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// inject dependency container
	c := container.NewContainer(db)

	// root route
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to POS Backend API 🚀"})
	})

	// API routes
	api := r.Group("/api")
	{
		// auth
		api.POST("/auth/signin", c.AuthHandler.SignIn)

		// user management (admin only)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.JWTAuth(), middleware.AdminOnly())
		{
			userRoutes.POST("/", c.UserHandler.CreateUser)
		}

		categories := api.Group("/categories")
		categories.Use(middleware.JWTAuth(), middleware.AdminOrStaffOnly())
		{
			categories.POST("", c.CategoryHandler.Create)
			categories.GET("", c.CategoryHandler.FindAll)
			categories.GET("/filter", c.CategoryHandler.FindWithFilter)
			categories.GET("/:id", c.CategoryHandler.FindByID)
			categories.PUT("/:id", c.CategoryHandler.Update)
			categories.DELETE("/:id", c.CategoryHandler.Delete)

		}
		product := api.Group("/products")
		product.Use(middleware.JWTAuth(), middleware.AdminOnly())
		{
			product.POST("", c.ProductHandler.Create)
			product.GET("", c.ProductHandler.FindAll)
			product.GET("/filter", c.ProductHandler.FindWithFilter)
			product.GET("/:id", c.ProductHandler.FindByID)
			product.PUT("/:id", c.ProductHandler.Update)
			product.DELETE("/:id", c.ProductHandler.Delete)
		}

		stockMovement := api.Group("/stock-movements")
		stockMovement.Use(middleware.JWTAuth(), middleware.AdminOrStaffOnly())
		{
			stockMovement.POST("", c.StockMovementHandler.Create)
			stockMovement.GET("", c.StockMovementHandler.FindAll)
			stockMovement.GET("/:id", c.StockMovementHandler.FindByIdProduct)
		}

		trx := api.Group("/transactions")
		trx.Use(middleware.JWTAuth(), middleware.AdminOrStaffOnly())
		{
			trx.POST("", c.TransactionHandler.Create)
		}

	}

	return r
}
