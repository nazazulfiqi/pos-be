package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(s service.StoreService) *StoreHandler {
	return &StoreHandler{service: s}
}

// Create godoc
// @Summary Create store
// @Description Create a new store. If owner_id is omitted, it's taken from the authenticated user.
// @Tags Store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateStoreRequest true "Create Store Request"
// @Success 201 {object} response.APIResponse{data=dto.StoreResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /stores/ [post]
func (h *StoreHandler) Create(ctx *gin.Context) {
	var req dto.CreateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}
	// if owner not provided, set from token user_id
	if req.OwnerID == nil {
		if uid, exists := ctx.Get("user_id"); exists {
			if us, ok := uid.(string); ok && us != "" {
				req.OwnerID = &us
			}
		}
	}
	st, err := h.service.Create(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(ctx, "Store created successfully", st)
}

// FindAll godoc
// @Summary Get all stores
// @Description Retrieve list of stores (optionally filtered by owner_id)
// @Tags Store
// @Produce json
// @Security BearerAuth
// @Param owner_id query string false "Owner ID"
// @Success 200 {object} response.APIResponse{data=[]dto.StoreResponse}
// @Failure 500 {object} response.APIResponse
// @Router /stores [get]
func (h *StoreHandler) FindAll(ctx *gin.Context) {
	var ownerID *string
	if v := ctx.Query("owner_id"); v != "" {
		ownerID = &v
	}
	stores, err := h.service.FindAll(ownerID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, "Stores fetched successfully", stores)
}

// FindByID godoc
// @Summary Get store by ID
// @Description Retrieve a single store by its ID
// @Tags Store
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 200 {object} response.APIResponse{data=dto.StoreResponse}
// @Failure 404 {object} response.APIResponse
// @Router /stores/{id} [get]
func (h *StoreHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	st, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Store fetched successfully", st)
}

// FindWithFilter godoc
// @Summary Get stores with filter and pagination
// @Description Retrieve stores using name filter, pagination and owner filter
// @Tags Store
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param owner_id query string false "Owner ID"
// @Success 200 {object} response.APIResponse{data=[]dto.StoreResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /stores/filter [get]
func (h *StoreHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.StoreFilter
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

	response.SuccessWithPagination(ctx, "Stores fetched successfully", data, meta)
}

// Update godoc
// @Summary Update store
// @Description Update an existing store by ID
// @Tags Store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Param request body dto.UpdateStoreRequest true "Update Store Request"
// @Success 200 {object} response.APIResponse{data=dto.StoreResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /stores/{id} [put]
func (h *StoreHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}
	// if owner not provided, set from token user_id
	if req.OwnerID == nil {
		if uid, exists := ctx.Get("user_id"); exists {
			if us, ok := uid.(string); ok && us != "" {
				req.OwnerID = &us
			}
		}
	}
	st, err := h.service.Update(id, req)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Store updated successfully", st)
}

// Delete godoc
// @Summary Delete store
// @Description Delete a store by ID
// @Tags Store
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /stores/{id} [delete]
func (h *StoreHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Store deleted successfully", nil)
}
