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
