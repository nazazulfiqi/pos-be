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

		permissions := api.Group("/permissions")
		permissions.Use(middleware.JWTAuth())
		{
			permissions.POST("", middleware.RequirePermission("permission.create"), c.PermissionHandler.Create)
			permissions.GET("", middleware.RequirePermission("permission.read"), c.PermissionHandler.FindAll)
			permissions.GET("/filter", middleware.RequirePermission("permission.read"), c.PermissionHandler.FindWithFilter)
			permissions.GET("/:id", middleware.RequirePermission("permission.read"), c.PermissionHandler.FindByID)
			permissions.PUT("/:id", middleware.RequirePermission("permission.update"), c.PermissionHandler.Update)
			permissions.DELETE("/:id", middleware.RequirePermission("permission.delete"), c.PermissionHandler.Delete)
		}

		roles := api.Group("/roles")
		roles.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			roles.POST("", middleware.RequirePermission("role.create"), c.RoleHandler.Create)
			roles.GET("", middleware.RequirePermission("role.read"), c.RoleHandler.FindAll)
			roles.GET("/filter", middleware.RequirePermission("role.read"), c.RoleHandler.FindWithFilter)
			roles.GET("/:id", middleware.RequirePermission("role.read"), c.RoleHandler.FindByID)
			roles.PUT("/:id", middleware.RequirePermission("role.update"), c.RoleHandler.Update)
			roles.DELETE("/:id", middleware.RequirePermission("role.delete"), c.RoleHandler.Delete)
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

		tenants := api.Group("/tenants")
		tenants.Use(middleware.JWTAuth(), middleware.TenantRequired())
		{
			tenants.GET("/filter", middleware.RequirePermission("tenant.read"), c.TenantHandler.FindWithFilter)
			tenants.POST("", middleware.RequirePermission("tenant.create"), c.TenantHandler.Create)
			tenants.GET("", middleware.RequirePermission("tenant.read"), c.TenantHandler.FindAll)
			tenants.GET("/:id", middleware.RequirePermission("tenant.read"), c.TenantHandler.FindByID)
			tenants.PUT("/:id", middleware.RequirePermission("tenant.update"), c.TenantHandler.Update)
			tenants.DELETE("/:id", middleware.RequirePermission("tenant.delete"), c.TenantHandler.Delete)
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
