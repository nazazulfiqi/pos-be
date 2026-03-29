package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(s service.TenantService) *TenantHandler {
	return &TenantHandler{service: s}
}

// Create godoc
// @Summary Create tenant
// @Description Create a new tenant. `store_id` must reference an existing store.
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTenantRequest true "Create Tenant Request"
// @Success 201 {object} response.APIResponse{data=dto.TenantResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /tenants/ [post]
func (h *TenantHandler) Create(ctx *gin.Context) {
	var req dto.CreateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}
	t, err := h.service.Create(req)
	if err != nil {
		s := err.Error()
		if strings.Contains(s, "violates foreign key") || strings.Contains(s, "23503") || strings.Contains(s, "fk_stores_tenants") {
			response.Error(ctx, http.StatusBadRequest, "store_id is invalid or does not exist")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(ctx, "Tenant created successfully", t)
}

// FindAll godoc
// @Summary Get all tenants
// @Description Retrieve list of tenants (optionally filtered by store_id)
// @Tags Tenant
// @Produce json
// @Security BearerAuth
// @Param store_id query string false "Store ID"
// @Success 200 {object} response.APIResponse{data=[]dto.TenantResponse}
// @Failure 500 {object} response.APIResponse
// @Router /tenants [get]
func (h *TenantHandler) FindAll(ctx *gin.Context) {
	var storeID *string
	if v := ctx.Query("store_id"); v != "" {
		storeID = &v
	}
	tenants, err := h.service.FindAll(storeID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, "Tenants fetched successfully", tenants)
}

// FindWithFilter godoc
// @Summary Get tenants with filter and pagination
// @Description Retrieve tenants using name filter, pagination and store filter
// @Tags Tenant
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param store_id query string false "Store ID"
// @Success 200 {object} response.APIResponse{data=[]dto.TenantResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /tenants/filter [get]
func (h *TenantHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.TenantFilter
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

	response.SuccessWithPagination(ctx, "Tenants fetched successfully", data, meta)
}

// FindByID godoc
// @Summary Get tenant by ID
// @Description Retrieve a single tenant by its ID
// @Tags Tenant
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.APIResponse{data=dto.TenantResponse}
// @Failure 404 {object} response.APIResponse
// @Router /tenants/{id} [get]
func (h *TenantHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	t, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Tenant fetched successfully", t)
}

// Update godoc
// @Summary Update tenant
// @Description Update an existing tenant by ID
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param request body dto.UpdateTenantRequest true "Update Tenant Request"
// @Success 200 {object} response.APIResponse{data=dto.TenantResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /tenants/{id} [put]
func (h *TenantHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}
	t, err := h.service.Update(id, req)
	if err != nil {
		s := err.Error()
		if s == "tenant not found" {
			response.Error(ctx, http.StatusNotFound, s)
			return
		}
		if strings.Contains(s, "violates foreign key") || strings.Contains(s, "23503") || strings.Contains(s, "fk_stores_tenants") {
			response.Error(ctx, http.StatusBadRequest, "store_id is invalid or does not exist")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, "Tenant updated successfully", t)
}

// Delete godoc
// @Summary Delete tenant
// @Description Delete a tenant by ID
// @Tags Tenant
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /tenants/{id} [delete]
func (h *TenantHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Tenant deleted successfully", nil)
}
