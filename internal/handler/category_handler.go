package handler

import (
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

// Create godoc
// @Summary Create category
// @Description Create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /categories/ [post]
func (h *CategoryHandler) Create(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	category, err := h.service.Create(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Category created successfully", category)
}

// FindAll godoc
// @Summary Get all categories
// @Description Retrieve list of categories (optionally filtered by tenant_id)
// @Tags Category
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]dto.CategoryResponse}
// @Failure 500 {object} response.APIResponse
// @Router /categories [get]
func (h *CategoryHandler) FindAll(ctx *gin.Context) {
	var tenantID *string
	if t := ctx.Query("tenant_id"); t != "" {
		tenantID = &t
	}

	categories, err := h.service.FindAll(tenantID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Categories fetched successfully", categories)

}

// FindWithFilter godoc
// @Summary Get categories with filter and pagination
// @Description Retrieve categories using search, pagination and tenant filter
// @Tags Category
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param tenant_id query string false "Tenant ID"
// @Success 200 {object} response.APIResponse{data=[]dto.CategoryResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /categories/filter [get]
func (h *CategoryHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.CategoryFilter
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

	response.SuccessWithPagination(ctx, "Categories fetched successfully", data, meta)
}

// FindByID godoc
// @Summary Get category by ID
// @Description Retrieve a single category by its ID
// @Tags Category
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure 404 {object} response.APIResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	category, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Category fetched successfully", category)
}

// Update godoc
// @Summary Update category
// @Description Update an existing category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.TranslateValidationError(err))
		return
	}

	category, err := h.service.Update(id, req)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, "Category updated successfully", category)
}

// Delete godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags Category
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.Success(ctx, "Category deleted successfully", nil)
}
