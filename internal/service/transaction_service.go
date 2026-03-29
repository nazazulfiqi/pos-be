package service

import (
	"bytes"
	"fmt"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"strconv"

	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userID string, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	PayTransaction(transactionID string, req dto.PayTransactionRequest) (*dto.PayTransactionResponse, error)
	GenerateReceiptPDF(transactionID string) ([]byte, error)
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

		if err := s.productRepo.DecreaseStock(trx, item.ProductID, item.Quantity); err != nil {
			trx.Rollback()
			return nil, err
		}

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
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found")
	}

	if transaction.PaymentStatus != model.PaymentStatusPending {
		return nil, fmt.Errorf("transaction is not pending, cannot process payment")
	}

	if transaction.PaymentMethod != model.PaymentMethodCash {
		return nil, fmt.Errorf("this endpoint only supports cash payment method")
	}

	if req.AmountTendered < transaction.TotalAmount {
		return nil, fmt.Errorf("insufficient amount: tendered %.2f but total is %.2f", req.AmountTendered, transaction.TotalAmount)
	}

	change := req.AmountTendered - transaction.TotalAmount

	if err := s.transactionRepo.UpdatePaymentStatus(s.db, transactionID, model.PaymentStatusSuccess); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %v", err)
	}

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

func (s *transactionService) GenerateReceiptPDF(transactionID string) ([]byte, error) {
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found")
	}

	if transaction.PaymentStatus != model.PaymentStatusSuccess {
		return nil, fmt.Errorf("receipt only available for paid transactions")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()
	pdf.SetMargins(15, 10, 15)

	drawLine := func(y float64) {
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetLineWidth(0.3)
		pdf.Line(15, y, 195, y)
	}

	// ─── Header ───
	tenantName := "POS Store"
	if transaction.Tenant != nil {
		tenantName = transaction.Tenant.Name
	}
	pdf.SetFont("Helvetica", "B", 18)
	pdf.CellFormat(0, 10, tenantName, "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.CellFormat(0, 5, "Jl. Contoh Alamat No. 123, Jakarta", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Telp: (021) 1234-5678", "", 1, "C", false, 0, "")

	lnY := pdf.GetY() + 2
	drawLine(lnY)
	pdf.SetY(lnY + 3)

	// ─── Title ───
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(0, 6, "STRUK / RECEIPT", "", 1, "C", false, 0, "")

	lnY = pdf.GetY() + 1
	drawLine(lnY)
	pdf.SetY(lnY + 3)

	// ─── Info ───
	pdf.SetFont("Helvetica", "", 9)
	infoLW := 35.0
	infoVW := 130.0

	pdf.CellFormat(infoLW, 5, "No. Transaksi", "", 0, "L", false, 0, "")
	pdf.CellFormat(infoVW, 5, ": "+transaction.ID, "", 1, "L", false, 0, "")

	cashierName := "-"
	if transaction.User.Name != "" {
		cashierName = transaction.User.Name
	}
	pdf.CellFormat(infoLW, 5, "Kasir", "", 0, "L", false, 0, "")
	pdf.CellFormat(infoVW, 5, ": "+cashierName, "", 1, "L", false, 0, "")

	if transaction.Customer != nil {
		pdf.CellFormat(infoLW, 5, "Customer", "", 0, "L", false, 0, "")
		pdf.CellFormat(infoVW, 5, ": "+transaction.Customer.Name, "", 1, "L", false, 0, "")
	}

	pdf.CellFormat(infoLW, 5, "Tanggal", "", 0, "L", false, 0, "")
	pdf.CellFormat(infoVW, 5, ": "+transaction.CreatedAt.Format("02/01/2006 15:04:05"), "", 1, "L", false, 0, "")

	lnY = pdf.GetY() + 2
	drawLine(lnY)
	pdf.SetY(lnY + 3)

	// ─── Table Header ───
	colItem := 75.0
	colQty := 20.0
	colPrice := 40.0
	colSub := 45.0

	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(colItem, 6, "Item", "", 0, "L", false, 0, "")
	pdf.CellFormat(colQty, 6, "Qty", "", 0, "C", false, 0, "")
	pdf.CellFormat(colPrice, 6, "Harga", "", 0, "R", false, 0, "")
	pdf.CellFormat(colSub, 6, "Subtotal", "", 1, "R", false, 0, "")

	lnY = pdf.GetY()
	drawLine(lnY)
	pdf.SetY(lnY + 2)

	// ─── Item Rows ───
	pdf.SetFont("Helvetica", "", 9)
	for _, item := range transaction.Items {
		productName := item.Product.Name
		if productName == "" {
			productName = item.ProductID
		}
		pdf.CellFormat(colItem, 5, productName, "", 0, "L", false, 0, "")
		pdf.CellFormat(colQty, 5, strconv.Itoa(item.Quantity), "", 0, "C", false, 0, "")
		pdf.CellFormat(colPrice, 5, fmt.Sprintf("Rp %.0f", item.Price), "", 0, "R", false, 0, "")
		pdf.CellFormat(colSub, 5, fmt.Sprintf("Rp %.0f", item.Subtotal), "", 1, "R", false, 0, "")
	}

	lnY = pdf.GetY() + 1
	drawLine(lnY)
	pdf.SetY(lnY + 3)

	// ─── Totals ───
	labelW := 95.0
	valW := 85.0

	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(labelW, 7, "TOTAL", "", 0, "L", false, 0, "")
	pdf.CellFormat(valW, 7, fmt.Sprintf("Rp %.0f", transaction.TotalAmount), "", 1, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	paymentMethodStr := "Cash"
	if transaction.PaymentMethod == model.PaymentMethodMidtrans {
		paymentMethodStr = "Midtrans"
	}
	pdf.CellFormat(labelW, 5, "Metode Pembayaran", "", 0, "L", false, 0, "")
	pdf.CellFormat(valW, 5, ": "+paymentMethodStr, "", 1, "R", false, 0, "")

	paymentStatusLabel := "Pending"
	switch transaction.PaymentStatus {
	case model.PaymentStatusSuccess:
		paymentStatusLabel = "LUNAS / PAID"
	case model.PaymentStatusCancel:
		paymentStatusLabel = "BATAL / CANCELLED"
	}
	pdf.CellFormat(labelW, 5, "Status", "", 0, "L", false, 0, "")
	pdf.CellFormat(valW, 5, ": "+paymentStatusLabel, "", 1, "R", false, 0, "")

	lnY = pdf.GetY() + 2
	drawLine(lnY)
	pdf.SetY(lnY + 4)

	// ─── Footer ───
	pdf.SetFont("Helvetica", "I", 9)
	pdf.CellFormat(0, 5, "Terima kasih atas kunjungan Anda!", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Thank you for your visit!", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Barang yang sudah dibeli tidak dapat ditukar atau dikembalikan.", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}
