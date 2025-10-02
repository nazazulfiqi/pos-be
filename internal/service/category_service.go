package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
)

type CategoryService interface {
	Create(req dto.CreateCategoryRequest) (dto.CategoryResponse, error)
	FindAll() ([]dto.CategoryResponse, error)
	FindWithFilter(filter dto.CategoryFilter) ([]dto.CategoryResponse, int64, error)
	FindByID(id uint) (dto.CategoryResponse, error)
	Update(id uint, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) Create(req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	category := model.Category{
		Name: req.Name,
	}

	if err := s.repo.Create(&category); err != nil {
		return dto.CategoryResponse{}, err
	}

	return dto.CategoryResponse{ID: category.ID, Name: category.Name}, nil
}

func (s *categoryService) FindAll() ([]dto.CategoryResponse, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponse, 0) // bukan var result []
	for _, c := range categories {
		result = append(result, dto.CategoryResponse{ID: c.ID, Name: c.Name})
	}
	return result, nil
}

func (s *categoryService) FindWithFilter(filter dto.CategoryFilter) ([]dto.CategoryResponse, int64, error) {
	categories, total, err := s.repo.FindWithFilter(filter.Search, filter.Page, filter.Limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CategoryResponse, 0) // bukan var result []
	for _, c := range categories {
		result = append(result, dto.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return result, total, nil
}

func (s *categoryService) FindByID(id uint) (dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("category not found")
	}
	return dto.CategoryResponse{ID: category.ID, Name: category.Name}, nil
}

func (s *categoryService) Update(id uint, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("category not found")
	}

	category.Name = req.Name

	if err := s.repo.Update(&category); err != nil {
		return dto.CategoryResponse{}, err
	}

	return dto.CategoryResponse{ID: category.ID, Name: category.Name}, nil
}

func (s *categoryService) Delete(id uint) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(&category)
}
