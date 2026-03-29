package handler

import (
	"mime/multipart"
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

// Create godoc
// @Summary Create product
// @Description Create a new product (multipart/form-data). Include `image` file optionally.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Product name"
// @Param sku formData string true "Product SKU"
// @Param category_id formData string true "Category ID"
// @Param tenant_id formData string true "Tenant ID"
// @Param price formData number true "Price"
// @Param stock formData int true "Stock"
// @Param image formData file false "Image file"
// @Success 201 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /products [post]
func (h *ProductHandler) Create(ctx *gin.Context) {
	var req dto.ProductCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// file upload (opsional)
	fileHeader, _ := ctx.FormFile("image")
	var file multipart.File
	var err error
	fileName := ""
	if fileHeader != nil {
		file, err = fileHeader.Open()
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "Failed to open uploaded file")
			return
		}
		defer file.Close()
		fileName = fileHeader.Filename
	}

	product, err := h.service.Create(req, file, fileName)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(ctx, "Product created successfully", product)
}

// Update godoc
// @Summary Update product
// @Description Update an existing product (multipart/form-data). `image` optional.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param name formData string true "Product name"
// @Param sku formData string true "Product SKU"
// @Param category_id formData string true "Category ID"
// @Param tenant_id formData string true "Tenant ID"
// @Param price formData number true "Price"
// @Param stock formData int true "Stock"
// @Param image formData file false "Image file"
// @Success 200 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /products/{id} [put]
func (h *ProductHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.ProductUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// file upload (opsional)
	fileHeader, _ := ctx.FormFile("image")
	var file multipart.File
	var err error
	fileName := ""
	if fileHeader != nil {
		file, err = fileHeader.Open()
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "Failed to open uploaded file")
			return
		}
		defer file.Close()
		fileName = fileHeader.Filename
	}

	product, err := h.service.Update(id, req, file, fileName)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Product updated successfully", product)
}

// Delete godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags Product
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Product deleted successfully", nil)
}

// FindByID godoc
// @Summary Get product by ID
// @Description Retrieve a single product by its ID
// @Tags Product
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure 404 {object} response.APIResponse
// @Router /products/{id} [get]
func (h *ProductHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := h.service.FindByID(id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Product not found")
		return
	}

	response.Success(ctx, "Product fetched successfully", product)
}

// FindAll godoc
// @Summary Get all products
// @Description Retrieve list of products (optionally filtered by tenant_id)
// @Tags Product
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "Tenant ID"
// @Success 200 {object} response.APIResponse{data=[]dto.ProductResponse}
// @Failure 500 {object} response.APIResponse
// @Router /products [get]
func (h *ProductHandler) FindAll(ctx *gin.Context) {
	var tenantID *string
	if t := ctx.Query("tenant_id"); t != "" {
		tenantID = &t
	}

	products, err := h.service.FindAll(tenantID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Products fetched successfully", products)
}

// FindWithFilter godoc
// @Summary Get products with filter and pagination
// @Description Retrieve products using filters (name, sku, category_id) and pagination
// @Tags Product
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param sku query string false "SKU"
// @Param category_id query string false "Category ID"
// @Param tenant_id query string false "Tenant ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.APIResponse{data=[]dto.ProductResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /products/filter [get]
func (h *ProductHandler) FindWithFilter(ctx *gin.Context) {
	var filter dto.ProductFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid query params")
		return
	}

	products, meta, err := h.service.FindWithFilter(filter)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPagination(ctx, "Products fetched successfully", products, meta)
}
