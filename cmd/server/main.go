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
