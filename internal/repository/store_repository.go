package repository

import (
	"pos-be/internal/model"

	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(store *model.Store) error
	FindAll(ownerID *string) ([]model.Store, error)
	FindWithFilter(name string, ownerID *string, page, limit int) ([]model.Store, int64, error)
	FindByID(id string) (model.Store, error)
	Update(store *model.Store) error
	Delete(store *model.Store) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db}
}

func (r *storeRepository) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepository) FindAll(ownerID *string) ([]model.Store, error) {
	var stores []model.Store
	query := r.db.Model(&model.Store{})
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}
	err := query.Find(&stores).Error
	return stores, err
}

func (r *storeRepository) FindWithFilter(name string, ownerID *string, page, limit int) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *storeRepository) FindByID(id string) (model.Store, error) {
	var store model.Store
	err := r.db.Where("id = ?", id).First(&store).Error
	return store, err
}

func (r *storeRepository) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(store *model.Store) error {
	return r.db.Delete(store).Error
}
