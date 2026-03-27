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
		// current user
		api.GET("/me", middleware.JWTAuth(), c.UserHandler.Me)
		// auth
		api.POST("/auth/signin", c.AuthHandler.SignIn)

		// user management (admin only)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			userRoutes.POST("/", middleware.RequirePermission("user.create"), c.UserHandler.CreateUser)
		}

		categories := api.Group("/categories")
		categories.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			categories.POST("", middleware.RequirePermission("category.create"), c.CategoryHandler.Create)
			categories.GET("", middleware.RequirePermission("category.read"), c.CategoryHandler.FindAll)
			categories.GET("/filter", middleware.RequirePermission("category.read"), c.CategoryHandler.FindWithFilter)
			categories.GET("/:id", middleware.RequirePermission("category.read"), c.CategoryHandler.FindByID)
			categories.PUT("/:id", middleware.RequirePermission("category.update"), c.CategoryHandler.Update)
			categories.DELETE("/:id", middleware.RequirePermission("category.delete"), c.CategoryHandler.Delete)

		}
		product := api.Group("/products")
		product.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			product.POST("", middleware.RequirePermission("product.create"), c.ProductHandler.Create)
			product.GET("", middleware.RequirePermission("product.read"), c.ProductHandler.FindAll)
			product.GET("/filter", middleware.RequirePermission("product.read"), c.ProductHandler.FindWithFilter)
			product.GET("/:id", middleware.RequirePermission("product.read"), c.ProductHandler.FindByID)
			product.PUT("/:id", middleware.RequirePermission("product.update"), c.ProductHandler.Update)
			product.DELETE("/:id", middleware.RequirePermission("product.delete"), c.ProductHandler.Delete)
		}

		stockMovement := api.Group("/stock-movements")
		stockMovement.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			stockMovement.POST("", middleware.RequirePermission("stockmovement.create"), c.StockMovementHandler.Create)
			stockMovement.GET("", middleware.RequirePermission("stockmovement.read"), c.StockMovementHandler.FindAll)
			stockMovement.GET("/:id", middleware.RequirePermission("stockmovement.read"), c.StockMovementHandler.FindByIdProduct)
		}

		trx := api.Group("/transactions")
		trx.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			trx.POST("", middleware.RequirePermission("transaction.create"), c.TransactionHandler.Create)
		}

	}

	stores := api.Group("/stores")
	stores.Use(middleware.JWTAuth(), middleware.TenantRequired())
	{
		stores.GET("/filter", middleware.RequirePermission("store.read"), c.StoreHandler.FindWithFilter)
		stores.POST("", middleware.RequirePermission("store.create"), c.StoreHandler.Create)
		stores.GET("", middleware.RequirePermission("store.read"), c.StoreHandler.FindAll)
		stores.GET("/:id", middleware.RequirePermission("store.read"), c.StoreHandler.FindByID)
		stores.PUT("/:id", middleware.RequirePermission("store.update"), c.StoreHandler.Update)
		stores.DELETE("/:id", middleware.RequirePermission("store.delete"), c.StoreHandler.Delete)
	}

	return r
}
