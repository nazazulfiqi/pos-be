package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type StockMovementHandler struct {
	service service.StockMovementService
}

func NewStockMovementHandler(service service.StockMovementService) *StockMovementHandler {
	return &StockMovementHandler{service}
}

// Create godoc
// @Summary Create stock movement
// @Description Create a new stock movement record
// @Tags StockMovement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.StockMovementCreateRequest true "Create Stock Movement Request"
// @Success 201 {object} response.APIResponse{data=dto.StockMovementResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /stock-movements/ [post]
func (h *StockMovementHandler) Create(ctx *gin.Context) {
	var req dto.StockMovementCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	movement, err := h.service.Create(req, nil)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Stock movement created successfully", movement)
}

// FindAll godoc
// @Summary Get all stock movements
// @Description Retrieve list of all stock movements
// @Tags StockMovement
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]dto.StockMovementResponse}
// @Failure 500 {object} response.APIResponse
// @Router /stock-movements [get]
func (h *StockMovementHandler) FindAll(ctx *gin.Context) {
	movements, err := h.service.FindAll()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Stock movements fetched successfully", movements)
}

// FindByProduct godoc
// @Summary Get stock movements by product
// @Description Retrieve stock movements for a specific product by its ID
// @Tags StockMovement
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.APIResponse{data=[]dto.StockMovementResponse}
// @Failure 500 {object} response.APIResponse
// @Router /stock-movements/product/{id} [get]
func (h *StockMovementHandler) FindByIdProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	movements, err := h.service.FindByProduct(productID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Stock movements fetched successfully", movements)
}
