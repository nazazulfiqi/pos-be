package dto

type TransactionItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	TenantID      string                   `json:"tenant_id" binding:"required"`
	CustomerID    *string                  `json:"customer_id,omitempty"`
	PaymentMethod uint                     `json:"payment_method" binding:"required,oneof=1 2"` // 1=Cash, 2=Midtrans
	Items         []TransactionItemRequest `json:"items" binding:"required,dive"`
}

type PayTransactionRequest struct {
	AmountTendered float64 `json:"amount_tendered" binding:"required,gt=0"`
}

type PayTransactionResponse struct {
	ID             string                    `json:"id"`
	TotalAmount    float64                   `json:"total_amount"`
	AmountTendered float64                   `json:"amount_tendered"`
	Change         float64                   `json:"change"`
	PaymentMethod  uint                      `json:"payment_method"`
	PaymentStatus  uint                      `json:"payment_status"`
	Items          []TransactionItemResponse `json:"items"`
}

type TransactionItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID            string                    `json:"id"`
	UserID        string                    `json:"user_id"`
	TenantID      string                    `json:"tenant_id"`
	CustomerID    *string                   `json:"customer_id,omitempty"`
	TotalAmount   float64                   `json:"total_amount"`
	PaymentStatus uint                      `json:"payment_status"`
	PaymentMethod uint                      `json:"payment_method,omitempty"` // 1=Cash, 2=Midtrans
	Items         []TransactionItemResponse `json:"items"`
	CreatedAt     string                    `json:"created_at"`
}
