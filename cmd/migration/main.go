package main

import (
	"log"

	"pos-be/config"
	"pos-be/internal/model"
)

func main() {
	// Init DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	// Auto migrate semua tabel sesuai blueprint
	err = db.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.Category{},
		&model.Product{},
		&model.StockMovement{},
		&model.Transaction{},
		&model.TransactionItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.WebhookLog{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Database migrated successfully")
}
