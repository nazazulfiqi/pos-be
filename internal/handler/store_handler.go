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

func (h *StoreHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	st, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Store fetched successfully", st)
}

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

func (h *StoreHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Store deleted successfully", nil)
}
