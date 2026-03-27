package dto

// --- CREATE ---
type StockMovementCreateRequest struct {
	ProductID     string  `json:"product_id" binding:"required"`
	Type          string  `json:"type" binding:"required,oneof=in out"` // validasi hanya "in" atau "out"
	Quantity      int     `json:"quantity" binding:"required,gt=0"`
	Note          string  `json:"note"`
	ReferenceID   string  `json:"reference_id"`
	ReferenceType string  `json:"reference_type" binding:"required"` // contoh: "transaction", "order", "purchase"
	TenantID      *string `json:"tenant_id,omitempty"`
}

// --- RESPONSE ---
type StockMovementResponse struct {
	ID            string `json:"id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	Type          string `json:"type"`
	Quantity      int    `json:"quantity"`
	Note          string `json:"note"`
	ReferenceID   string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type"`
	CreatedAt     string `json:"created_at"`
}
