package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

// SignIn godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SignInRequest true "Sign In Request"
// @Success 200 {object} response.APIResponse{data=dto.SignInResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/signin [post]
func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var req dto.SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := h.service.SignIn(req)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(ctx, "Login successful", res)
}
