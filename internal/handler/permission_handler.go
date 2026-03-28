package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{service}
}

func (h *PermissionHandler) Create(ctx *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	permission, err := h.service.Create(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Permission created successfully", permission)
}

func (h *PermissionHandler) FindAll(ctx *gin.Context) {
	permissions, err := h.service.FindAll()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Permissions fetched successfully", permissions)
}

func (h *PermissionHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.PermissionFilter
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

	response.SuccessWithPagination(ctx, "Permissions fetched successfully", data, meta)
}

func (h *PermissionHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	permission, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Permission fetched successfully", permission)
}

func (h *PermissionHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	permission, err := h.service.Update(id, req)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, "Permission updated successfully", permission)
}

func (h *PermissionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Permission deleted successfully", nil)
}