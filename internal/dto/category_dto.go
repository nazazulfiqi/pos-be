package dto

type CreateCategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	TenantID *string `json:"tenant_id" binding:"required"` // Menambahkan tenant_id sebagai field yang wajib diisi
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryFilter struct {
	Search   string  `form:"search"`
	Page     int     `form:"page" binding:"min=1"`
	Limit    int     `form:"limit" binding:"min=1,max=100"`
	TenantID *string `form:"tenant_id"`
}

// (Tenant-aware CreateCategoryRequest is above)
