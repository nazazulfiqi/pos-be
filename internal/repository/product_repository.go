package repository

import (
	"fmt"
	"pos-be/internal/dto"
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id string) error
	FindByID(id string) (*model.Product, error)
	FindAll(tenantID *string) ([]model.Product, error)
	FindWithFilter(filter dto.ProductFilter) ([]model.Product, int64, error)
	DecreaseStock(tx *gorm.DB, productID string, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(tenantID *string) ([]model.Product, error) {
	var products []model.Product
	query := r.db.Preload("Category").Model(&model.Product{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindWithFilter(filter dto.ProductFilter) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Preload("Category")

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.SKU != "" {
		query = query.Where("sku ILIKE ?", "%"+filter.SKU+"%")
	}

	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}

	// total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) DecreaseStock(tx *gorm.DB, productID string, quantity int) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	// Kurangi stok, pastikan stok tidak kurang dari 0
	result := db.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock for product ID %s", productID)
	}

	return nil
}
