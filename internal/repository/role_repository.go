package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *model.Role) error
	FindAll(tenantID *string) ([]model.Role, error)
	FindWithFilter(search string, page, limit int, tenantID, storeID *string) ([]model.Role, int64, error)
	FindByID(id string) (model.Role, error)
	FindBySlug(slug string) (model.Role, error)
	Update(role *model.Role) error
	Delete(role *model.Role) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindAll(tenantID *string) ([]model.Role, error) {
	var roles []model.Role
	query := r.db.Model(&model.Role{}).Preload("Permissions")
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	err := query.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindWithFilter(search string, page, limit int, tenantID, storeID *string) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.Model(&model.Role{}).Preload("Permissions")

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination offset
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) FindByID(id string) (model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").Where("id = ?", id).First(&role).Error
	return role, err
}

func (r *roleRepository) FindBySlug(slug string) (model.Role, error) {
	var role model.Role
	err := r.db.Where("slug = ?", slug).First(&role).Error
	return role, err
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(role *model.Role) error {
	return r.db.Delete(role).Error
}