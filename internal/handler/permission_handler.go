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

// Create godoc
// @Summary Create permission
// @Description Create a new permission
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePermissionRequest true "Create Permission Request"
// @Success 201 {object} response.APIResponse{data=dto.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /permissions/ [post]
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

// FindAll godoc
// @Summary Get all permissions
// @Description Retrieve list of permissions
// @Tags Permission
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]dto.PermissionResponse}
// @Failure 500 {object} response.APIResponse
// @Router /permissions [get]
func (h *PermissionHandler) FindAll(ctx *gin.Context) {
	permissions, err := h.service.FindAll()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Permissions fetched successfully", permissions)
}

// FindWithFilter godoc
// @Summary Get permissions with filter and pagination
// @Description Retrieve permissions using name filter and pagination
// @Tags Permission
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.APIResponse{data=[]dto.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /permissions/filter [get]
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

// FindByID godoc
// @Summary Get permission by ID
// @Description Retrieve a single permission by its ID
// @Tags Permission
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID"
// @Success 200 {object} response.APIResponse{data=dto.PermissionResponse}
// @Failure 404 {object} response.APIResponse
// @Router /permissions/{id} [get]
func (h *PermissionHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	permission, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Permission fetched successfully", permission)
}

// Update godoc
// @Summary Update permission
// @Description Update an existing permission by ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID"
// @Param request body dto.UpdatePermissionRequest true "Update Permission Request"
// @Success 200 {object} response.APIResponse{data=dto.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /permissions/{id} [put]
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

// Delete godoc
// @Summary Delete permission
// @Description Delete a permission by ID
// @Tags Permission
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /permissions/{id} [delete]
func (h *PermissionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Permission deleted successfully", nil)
}
