package dto

// Request
type CreateUserRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=5"`
	RoleIDs  []string `json:"role_ids" binding:"required"`
	TenantID *string  `json:"tenant_id"`
}

// Response
type UserResponse struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Roles    []string `json:"role_ids"`
	TenantID *string  `json:"tenant_id,omitempty"`
}

type TenantDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	StoreID *string `json:"store_id,omitempty"`
}

type StoreDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	OwnerID *string `json:"owner_id,omitempty"`
}

type MeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions"`
	Tenant      *TenantDTO `json:"tenant,omitempty"`
	Store       *StoreDTO  `json:"store,omitempty"`
}
