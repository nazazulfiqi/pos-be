package dto

type CreateTenantRequest struct {
	Name    string  `json:"name" binding:"required"`
	StoreID *string `json:"store_id" binding:"required"`
}

type UpdateTenantRequest struct {
	Name    string  `json:"name" binding:"required"`
	StoreID *string `json:"store_id" binding:"required"`
}

type TenantResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	StoreID *string `json:"store_id,omitempty"`
}

type TenantFilter struct {
	Name    string  `form:"name"`
	Page    int     `form:"page" binding:"min=1"`
	Limit   int     `form:"limit" binding:"min=1,max=100"`
	StoreID *string `form:"store_id"`
}
