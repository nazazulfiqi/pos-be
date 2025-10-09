package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(tx *gorm.DB, movement *model.StockMovement) error
	FindAll() ([]model.StockMovement, error)
	FindByProduct(productID uint) ([]model.StockMovement, error)
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db}
}

func (r *stockMovementRepository) Create(tx *gorm.DB, movement *model.StockMovement) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(movement).Error
}

func (r *stockMovementRepository) FindAll() ([]model.StockMovement, error) {
	var movements []model.StockMovement
	err := r.db.Preload("Product").Order("created_at desc").Find(&movements).Error
	return movements, err
}

func (r *stockMovementRepository) FindByProduct(productID uint) ([]model.StockMovement, error) {
	var movements []model.StockMovement
	err := r.db.Preload("Product").
		Where("product_id = ?", productID).
		Order("created_at desc").
		Find(&movements).Error
	return movements, err
}
