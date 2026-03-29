// @title POS Backend API
// @version 1.0
// @description REST API for Point of Sale (POS) application with multi-tenant support.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

//go:generate swag init -g cmd/server/main.go --parseInternal -d .

package main

import (
	"log"
	"os"

	"pos-be/config"
	"pos-be/internal/router"
)

func main() {
	// Init DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	// Setup router dengan DB
	r := router.SetupRouter(db)

	// Ambil port dari env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 POS Backend running on port", port)
	r.Run(":" + port)
}
