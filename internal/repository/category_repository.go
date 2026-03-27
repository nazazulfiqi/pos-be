package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindAll(tenantID *string) ([]model.Category, error)
	FindWithFilter(search string, page, limit int, tenantID *string) ([]model.Category, int64, error)
	FindByID(id string) (model.Category, error)
	Update(category *model.Category) error
	Delete(category *model.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAll(tenantID *string) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.Model(&model.Category{})
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindWithFilter(search string, page, limit int, tenantID *string) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.Model(&model.Category{})

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination offset
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(id string) (model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return category, err
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(category *model.Category) error {
	return r.db.Delete(category).Error
}
