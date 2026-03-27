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

func (h *TenantHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	t, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Tenant fetched successfully", t)
}

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

func (h *TenantHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Tenant deleted successfully", nil)
}
