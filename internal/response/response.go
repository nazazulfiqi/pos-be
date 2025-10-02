package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type PaginationMeta struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
}

type PaginatedResponse struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Meta       PaginationMeta `json:"meta"`
}

func SuccessWithPagination(ctx *gin.Context, message string, data interface{}, meta PaginationMeta) {
	ctx.JSON(http.StatusOK, PaginatedResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
		Meta:       meta,
	})
}

func Success(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, APIResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	})
}

func Created(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, APIResponse{
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	})
}

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
	})
}
