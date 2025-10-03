package dto

// --- CREATE ---
type ProductCreateRequest struct {
	Name       string  `form:"name" binding:"required"`
	SKU        string  `form:"sku" binding:"required"`
	CategoryID uint    `form:"category_id" binding:"required"`
	Price      float64 `form:"price" binding:"required"`
	Stock      int     `form:"stock" binding:"required"`
	// image di-handle di handler pakai ctx.FormFile("image")
}

// --- UPDATE ---
type ProductUpdateRequest struct {
	Name       string  `form:"name" binding:"required"`
	SKU        string  `form:"sku" binding:"required"`
	CategoryID uint    `form:"category_id" binding:"required"`
	Price      float64 `form:"price" binding:"required"`
	Stock      int     `form:"stock" binding:"required"`
	// image optional → kalau tidak ada, pakai gambar lama
}

// --- RESPONSE ---
type ProductResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	SKU       string  `json:"sku"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	ImageURL  string  `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	PublicID  string  `json:"-"`
}

// --- FILTER ---
type ProductFilter struct {
	Name       string `form:"name"`
	SKU        string `form:"sku"`
	CategoryID uint   `form:"category_id"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}
