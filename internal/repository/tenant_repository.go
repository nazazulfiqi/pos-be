package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type TenantRepository interface {
	Create(tenant *model.Tenant) error
	FindAll(storeID *string) ([]model.Tenant, error)
	FindWithFilter(name string, storeID *string, page, limit int) ([]model.Tenant, int64, error)
	FindByID(id string) (model.Tenant, error)
	Update(tenant *model.Tenant) error
	Delete(tenant *model.Tenant) error
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db}
}

func (r *tenantRepository) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepository) FindAll(storeID *string) ([]model.Tenant, error) {
	var tenants []model.Tenant
	query := r.db.Model(&model.Tenant{})
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	err := query.Find(&tenants).Error
	return tenants, err
}

func (r *tenantRepository) FindWithFilter(name string, storeID *string, page, limit int) ([]model.Tenant, int64, error) {
	var tenants []model.Tenant
	var total int64

	query := r.db.Model(&model.Tenant{})
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

func (r *tenantRepository) FindByID(id string) (model.Tenant, error) {
	var t model.Tenant
	err := r.db.Where("id = ?", id).First(&t).Error
	return t, err
}

func (r *tenantRepository) Update(tenant *model.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *tenantRepository) Delete(tenant *model.Tenant) error {
	return r.db.Delete(tenant).Error
}
