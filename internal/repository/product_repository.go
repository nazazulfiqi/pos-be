package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uint) error
	FindByID(id uint) (*model.Product, error)
	FindAll() ([]model.Product, error)
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

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
