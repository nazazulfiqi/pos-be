package repository

import (
	"fmt"
	"pos-be/internal/model"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GenerateTransactionID() (string, error)
	CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error
	FindByID(id string) (*model.Transaction, error)
	UpdatePaymentStatus(tx *gorm.DB, id string, status uint) error
	WithTrx(trx *gorm.DB) TransactionRepository
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) WithTrx(trx *gorm.DB) TransactionRepository {
	if trx == nil {
		return r
	}
	return &transactionRepository{db: trx}
}

// generate format TRX-25-0001
func (r *transactionRepository) GenerateTransactionID() (string, error) {
	year := time.Now().Year() % 100
	prefix := fmt.Sprintf("TRX-%02d-", year)

	var lastTrx model.Transaction
	if err := r.db.
		Where("id LIKE ?", prefix+"%").
		Order("id DESC").
		First(&lastTrx).Error; err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	lastNumber := 0
	if lastTrx.ID != "" {
		fmt.Sscanf(lastTrx.ID, prefix+"%04d", &lastNumber)
	}
	newNumber := lastNumber + 1
	newID := fmt.Sprintf("%s%04d", prefix, newNumber)

	return newID, nil
}

func (r *transactionRepository) CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.Preload("Items.Product").Preload("Tenant").Preload("User").Preload("Customer").Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdatePaymentStatus(tx *gorm.DB, id string, status uint) error {
	return tx.Model(&model.Transaction{}).Where("id = ?", id).Update("payment_status", status).Error
}
