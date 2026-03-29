package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with role assignment
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Create User Request"
// @Success 201 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Router /users/ [post]
func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// gunakan helper translate error
		message := response.TranslateValidationError(err)
		response.Error(ctx, http.StatusBadRequest, message)
		return
	}

	userRes, err := h.service.CreateUser(req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(ctx, "User created successfully", userRes)
}

// Me godoc
// @Summary Get current user
// @Description Retrieve profile of the authenticated user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 401 {object} response.APIResponse
// @Router /me [get]
func (h *UserHandler) Me(ctx *gin.Context) {
	uid, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, ok := uid.(string)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Invalid token user id")
		return
	}

	me, err := h.service.GetMe(userID)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, "OK", me)
}
