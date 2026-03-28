package dto

import "time"

type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   string   `json:"description"`
	TenantID      *string  `json:"tenant_id"`
	StoreID       *string  `json:"store_id"`
	PermissionIDs []string `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

type RoleResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Slug         string               `json:"slug"`
	Description  string               `json:"description"`
	TenantID     *string              `json:"tenant_id"`
	StoreID      *string              `json:"store_id"`
	Permissions  []PermissionResponse `json:"permissions"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type RoleFilter struct {
	Search   string  `form:"search"`
	Page     int     `form:"page" binding:"min=1"`
	Limit    int     `form:"limit" binding:"min=1,max=100"`
	TenantID *string `form:"tenant_id"`
	StoreID  *string `form:"store_id"`
}