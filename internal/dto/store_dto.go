package dto

type CreateStoreRequest struct {
	Name    string  `json:"name" binding:"required"`
	OwnerID *string `json:"owner_id"`
}

type UpdateStoreRequest struct {
	Name    string  `json:"name" binding:"required"`
	OwnerID *string `json:"owner_id"`
}

type StoreResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	OwnerID *string `json:"owner_id,omitempty"`
}

type StoreFilter struct {
	Name    string  `form:"name"`
	Page    int     `form:"page" binding:"min=1"`
	Limit   int     `form:"limit" binding:"min=1,max=100"`
	OwnerID *string `form:"owner_id"`
}
