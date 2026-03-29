package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Create godoc
// @Summary Create transaction
// @Description Create a new transaction. `tenant_id` required; `user_id` is taken from the authenticated user.
// @Tags Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Create Transaction Request"
// @Success 201 {object} response.APIResponse{data=dto.TransactionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /transactions/ [post]
func (h *TransactionHandler) Create(ctx *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user_id from JWT token
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	result, err := h.service.CreateTransaction(userID.(string), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Transaction created successfully", result)
}

// Pay godoc
// @Summary Pay transaction
// @Description Process payment for a transaction by ID
// @Tags Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Param request body dto.PayTransactionRequest true "Payment Request"
// @Success 200 {object} response.APIResponse{data=dto.PayTransactionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /transactions/{id}/pay [post]
func (h *TransactionHandler) Pay(ctx *gin.Context) {
	transactionID := ctx.Param("id")

	var req dto.PayTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := h.service.PayTransaction(transactionID, req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, "Payment processed successfully", result)
}

// Receipt godoc
// @Summary Get transaction receipt PDF
// @Description Generate and download PDF receipt for a transaction
// @Tags Transaction
// @Produce application/pdf
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {file} file "PDF receipt"
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /transactions/{id}/receipt [get]
func (h *TransactionHandler) Receipt(ctx *gin.Context) {
	transactionID := ctx.Param("id")

	pdfBytes, err := h.service.GenerateReceiptPDF(transactionID)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=receipt_"+transactionID+".pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
