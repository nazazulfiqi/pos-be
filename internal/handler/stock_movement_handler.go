package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StockMovementHandler struct {
	service service.StockMovementService
}

func NewStockMovementHandler(service service.StockMovementService) *StockMovementHandler {
	return &StockMovementHandler{service}
}

func (h *StockMovementHandler) Create(ctx *gin.Context) {
	var req dto.StockMovementCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	movement, err := h.service.Create(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Stock movement created successfully", movement)
}

func (h *StockMovementHandler) FindAll(ctx *gin.Context) {
	movements, err := h.service.FindAll()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Stock movements fetched successfully", movements)
}

func (h *StockMovementHandler) FindByProduct(ctx *gin.Context) {
	productIDParam := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid product_id")
		return
	}

	movements, err := h.service.FindByProduct(uint(productID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Stock movements fetched successfully", movements)
}
