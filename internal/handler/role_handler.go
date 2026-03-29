package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service service.RoleService
}

func NewRoleHandler(service service.RoleService) *RoleHandler {
	return &RoleHandler{service}
}

// Create godoc
// @Summary Create role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRoleRequest true "Create Role Request"
// @Success 201 {object} response.APIResponse{data=dto.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /roles/ [post]
func (h *RoleHandler) Create(ctx *gin.Context) {
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	role, err := h.service.Create(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Role created successfully", role)
}

// FindAll godoc
// @Summary Get all roles
// @Description Retrieve list of roles (optionally filtered by tenant_id)
// @Tags Role
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "Tenant ID"
// @Success 200 {object} response.APIResponse{data=[]dto.RoleResponse}
// @Failure 500 {object} response.APIResponse
// @Router /roles [get]
func (h *RoleHandler) FindAll(ctx *gin.Context) {
	var tenantID *string
	if t := ctx.Query("tenant_id"); t != "" {
		tenantID = &t
	}

	roles, err := h.service.FindAll(tenantID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Roles fetched successfully", roles)
}

// FindWithFilter godoc
// @Summary Get roles with filter and pagination
// @Description Retrieve roles using name filter, pagination and tenant filter
// @Tags Role
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param tenant_id query string false "Tenant ID"
// @Success 200 {object} response.APIResponse{data=[]dto.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /roles/filter [get]
func (h *RoleHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.RoleFilter
	// default pagination
	filter.Page = 1
	filter.Limit = 10

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid query params")
		return
	}

	data, meta, err := h.service.FindWithFilter(filter)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPagination(ctx, "Roles fetched successfully", data, meta)
}

// FindByID godoc
// @Summary Get role by ID
// @Description Retrieve a single role by its ID
// @Tags Role
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} response.APIResponse{data=dto.RoleResponse}
// @Failure 404 {object} response.APIResponse
// @Router /roles/{id} [get]
func (h *RoleHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Role fetched successfully", role)
}

// Update godoc
// @Summary Update role
// @Description Update an existing role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Update Role Request"
// @Success 200 {object} response.APIResponse{data=dto.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	role, err := h.service.Update(id, req)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, "Role updated successfully", role)
}

// Delete godoc
// @Summary Delete role
// @Description Delete a role by ID
// @Tags Role
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Role deleted successfully", nil)
}
