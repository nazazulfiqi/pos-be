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
