package handler

import (
	"mime/multipart"
	"net/http"
	"pos-be/internal/dto"
	"pos-be/internal/response"
	"pos-be/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

// --- CREATE ---
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

// --- UPDATE ---
func (h *ProductHandler) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req dto.ProductUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// file upload (opsional)
	fileHeader, _ := ctx.FormFile("image")
	var file multipart.File
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

	product, err := h.service.Update(uint(id), req, file, fileName)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Product updated successfully", product)
}

// --- DELETE ---
func (h *ProductHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Product deleted successfully", nil)
}

// --- FIND BY ID ---
func (h *ProductHandler) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.FindByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Product not found")
		return
	}

	response.Success(ctx, "Product fetched successfully", product)
}

// --- FIND ALL ---
func (h *ProductHandler) FindAll(ctx *gin.Context) {
	products, err := h.service.FindAll()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, "Products fetched successfully", products)
}

// --- FIND WITH FILTER ---
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
