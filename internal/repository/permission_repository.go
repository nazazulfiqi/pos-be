package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *model.Permission) error
	FindAll() ([]model.Permission, error)
	FindWithFilter(search string, page, limit int) ([]model.Permission, int64, error)
	FindByID(id string) (model.Permission, error)
	FindBySlug(slug string) (model.Permission, error)
	Update(permission *model.Permission) error
	Delete(permission *model.Permission) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) FindAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Model(&model.Permission{}).Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindWithFilter(search string, page, limit int) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	query := r.db.Model(&model.Permission{})

	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination offset
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *permissionRepository) FindByID(id string) (model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	return permission, err
}

func (r *permissionRepository) FindBySlug(slug string) (model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("slug = ?", slug).First(&permission).Error
	return permission, err
}

func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(permission *model.Permission) error {
	return r.db.Delete(permission).Error
}