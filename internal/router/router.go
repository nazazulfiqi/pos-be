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
	}

	return r
}
