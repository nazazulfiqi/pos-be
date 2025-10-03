package service

import (
	"mime/multipart"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"pos-be/internal/response"
	"pos-be/utils"
	"time"
)

type ProductService interface {
	Create(req dto.ProductCreateRequest, file multipart.File, fileName string) (dto.ProductResponse, error)
	Update(id uint, req dto.ProductUpdateRequest, file multipart.File, fileName string) (dto.ProductResponse, error)
	Delete(id uint) error
	FindByID(id uint) (dto.ProductResponse, error)
	FindAll() ([]dto.ProductResponse, error)
	FindWithFilter(filter dto.ProductFilter) ([]dto.ProductResponse, response.PaginationMeta, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

// --- CREATE ---
func (s *productService) Create(req dto.ProductCreateRequest, file multipart.File, fileName string) (dto.ProductResponse, error) {
	imageURL := ""
	publicID := ""

	// kalau ada file upload
	if file != nil {
		uploadedURL, uploadedID, err := utils.UploadToCloudinary(file, fileName)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		imageURL = uploadedURL
		publicID = uploadedID
	}

	product := model.Product{
		Name:       req.Name,
		SKU:        req.SKU,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		Stock:      req.Stock,
		ImageURL:   imageURL,
		PublicID:   publicID,
	}

	if err := s.repo.Create(&product); err != nil {
		return dto.ProductResponse{}, err
	}

	return s.toResponse(product), nil
}

// --- UPDATE ---
func (s *productService) Update(id uint, req dto.ProductUpdateRequest, file multipart.File, fileName string) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// kalau ada file baru → hapus lama + upload baru
	if file != nil {
		if product.PublicID != "" {
			_ = utils.DeleteFromCloudinary(product.PublicID)
		}
		uploadedURL, uploadedID, err := utils.UploadToCloudinary(file, fileName)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		product.ImageURL = uploadedURL
		product.PublicID = uploadedID
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.CategoryID = req.CategoryID
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.repo.Update(product); err != nil {
		return dto.ProductResponse{}, err
	}

	return s.toResponse(*product), nil
}

// --- DELETE ---
func (s *productService) Delete(id uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if product.PublicID != "" {
		_ = utils.DeleteFromCloudinary(product.PublicID)
	}

	return s.repo.Delete(id)
}

// --- FIND BY ID ---
func (s *productService) FindByID(id uint) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return s.toResponse(*product), nil
}

// --- FIND ALL ---
func (s *productService) FindAll() ([]dto.ProductResponse, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.ProductResponse, 0)
	for _, p := range products {
		res = append(res, s.toResponse(p))
	}
	return res, nil
}

// --- Helper mapping ---
func (s *productService) toResponse(p model.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		SKU:       p.SKU,
		Category:  p.Category.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		ImageURL:  p.ImageURL,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		PublicID:  p.PublicID,
	}
}

func (s *productService) FindWithFilter(filter dto.ProductFilter) ([]dto.ProductResponse, response.PaginationMeta, error) {
	products, total, err := s.repo.FindWithFilter(filter)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	// mapping ke response DTO
	result := []dto.ProductResponse{}
	for _, p := range products {
		result = append(result, dto.ProductResponse{
			ID:        p.ID,
			Name:      p.Name,
			SKU:       p.SKU,
			Category:  p.Category.Name,
			Price:     p.Price,
			Stock:     p.Stock,
			ImageURL:  p.ImageURL,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
			UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
			PublicID:  p.PublicID,
		})
	}

	if result == nil {
		result = []dto.ProductResponse{}
	}

	// pagination meta dihitung di service, bukan handler
	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	meta := response.PaginationMeta{
		TotalRecords: total,
		TotalPages:   totalPages,
		CurrentPage:  filter.Page,
		PageSize:     filter.Limit,
	}

	return result, meta, nil
}
