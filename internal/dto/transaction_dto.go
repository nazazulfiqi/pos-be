package dto

type TransactionItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	UserID        string                   `json:"user_id" binding:"required"`
	CustomerID    *string                  `json:"customer_id,omitempty"`
	TenantID      *string                  `json:"tenant_id,omitempty"`
	PaymentMethod string                   `json:"payment_method" binding:"required"`
	Items         []TransactionItemRequest `json:"items" binding:"required,dive"`
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
	CustomerID    *string                   `json:"customer_id,omitempty"`
	TotalAmount   float64                   `json:"total_amount"`
	PaymentMethod string                    `json:"payment_method"`
	Items         []TransactionItemResponse `json:"items"`
	CreatedAt     string                    `json:"created_at"`
}
