package service

import (
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
}

type transactionService struct {
	db                *gorm.DB
	transactionRepo   repository.TransactionRepository
	productRepo       repository.ProductRepository
	stockMovementRepo repository.StockMovementRepository
}

func NewTransactionService(
	db *gorm.DB,
	transactionRepo repository.TransactionRepository,
	productRepo repository.ProductRepository,
	stockMovementRepo repository.StockMovementRepository,
) TransactionService {
	return &transactionService{
		db:                db,
		transactionRepo:   transactionRepo,
		productRepo:       productRepo,
		stockMovementRepo: stockMovementRepo,
	}
}

func (s *transactionService) CreateTransaction(req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	trx := s.db.Begin()

	id, err := s.transactionRepo.GenerateTransactionID()
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	var total float64
	var items []model.TransactionItem

	for _, item := range req.Items {
		subtotal := float64(item.Quantity) * item.Price
		total += subtotal

		items = append(items, model.TransactionItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  subtotal,
		})

		// kurangi stok
		if err := s.productRepo.DecreaseStock(trx, item.ProductID, item.Quantity); err != nil {
			trx.Rollback()
			return nil, err
		}

		// buat record stock movement (out)
		sm := model.StockMovement{
			ProductID:     item.ProductID,
			Type:          "out",
			Quantity:      item.Quantity,
			Note:          "Transaction sale",
			ReferenceID:   id,
			ReferenceType: "transaction",
		}
		if err := s.stockMovementRepo.Create(trx, &sm); err != nil {
			trx.Rollback()
			return nil, err
		}
	}

	transaction := &model.Transaction{
		ID:            id,
		UserID:        req.UserID,
		CustomerID:    req.CustomerID,
		PaymentMethod: req.PaymentMethod,
		TotalAmount:   total,
		Items:         items,
	}

	if err := s.transactionRepo.CreateTransaction(trx, transaction); err != nil {
		trx.Rollback()
		return nil, err
	}

	trx.Commit()

	resp := dto.TransactionResponse{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		CustomerID:    transaction.CustomerID,
		PaymentMethod: transaction.PaymentMethod,
		TotalAmount:   transaction.TotalAmount,
		CreatedAt:     transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	for _, i := range transaction.Items {
		resp.Items = append(resp.Items, dto.TransactionItemResponse{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
			Price:     i.Price,
			Subtotal:  i.Subtotal,
		})
	}

	return &resp, nil
}
