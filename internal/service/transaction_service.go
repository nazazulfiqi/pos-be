package service

import (
	"fmt"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	PayTransaction(transactionID string, req dto.PayTransactionRequest) (*dto.PayTransactionResponse, error)
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

func (s *transactionService) CreateTransaction(userID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	trx := s.db.Begin()

	id, err := s.transactionRepo.GenerateTransactionID()
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	var total float64
	var items []model.TransactionItem

	for _, item := range req.Items {
		// Lookup product to get current price
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			trx.Rollback()
			return nil, fmt.Errorf("product not found: %s", item.ProductID)
		}

		subtotal := float64(item.Quantity) * product.Price
		total += subtotal

		items = append(items, model.TransactionItem{
			ProductID: item.ProductID,
			TenantID:  req.TenantID,
			Quantity:  item.Quantity,
			Price:     product.Price,
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
			TenantID:      &req.TenantID,
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
		UserID:        userID,
		CustomerID:    req.CustomerID,
		TenantID:      req.TenantID,
		TotalAmount:   total,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: model.PaymentStatusPending,
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
		TenantID:      transaction.TenantID,
		CustomerID:    transaction.CustomerID,
		TotalAmount:   transaction.TotalAmount,
		PaymentStatus: transaction.PaymentStatus,
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

func (s *transactionService) PayTransaction(transactionID string, req dto.PayTransactionRequest) (*dto.PayTransactionResponse, error) {
	// Find transaction
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found")
	}

	// Validate payment status must be pending
	if transaction.PaymentStatus != model.PaymentStatusPending {
		return nil, fmt.Errorf("transaction is not pending, cannot process payment")
	}

	// Validate payment method must be cash (1)
	if transaction.PaymentMethod != model.PaymentMethodCash {
		return nil, fmt.Errorf("this endpoint only supports cash payment method")
	}

	// Validate amount tendered must be >= total amount
	if req.AmountTendered < transaction.TotalAmount {
		return nil, fmt.Errorf("insufficient amount: tendered %.2f but total is %.2f", req.AmountTendered, transaction.TotalAmount)
	}

	// Calculate change
	change := req.AmountTendered - transaction.TotalAmount

	// Update payment status to success
	if err := s.transactionRepo.UpdatePaymentStatus(s.db, transactionID, model.PaymentStatusSuccess); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %v", err)
	}

	// Build response
	resp := &dto.PayTransactionResponse{
		ID:             transaction.ID,
		TotalAmount:    transaction.TotalAmount,
		AmountTendered: req.AmountTendered,
		Change:         change,
		PaymentMethod:  transaction.PaymentMethod,
		PaymentStatus:  model.PaymentStatusSuccess,
		Items:          make([]dto.TransactionItemResponse, 0, len(transaction.Items)),
	}

	for _, item := range transaction.Items {
		resp.Items = append(resp.Items, dto.TransactionItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
		})
	}

	return resp, nil
}
