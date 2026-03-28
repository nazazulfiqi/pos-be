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

func (h *RoleHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Role fetched successfully", role)
}

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

func (h *RoleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Role deleted successfully", nil)
}