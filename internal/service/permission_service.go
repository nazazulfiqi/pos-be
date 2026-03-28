package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"pos-be/internal/response"
)

type PermissionService interface {
	Create(req dto.CreatePermissionRequest) (dto.PermissionResponse, error)
	FindAll() ([]dto.PermissionResponse, error)
	FindWithFilter(filter dto.PermissionFilter) ([]dto.PermissionResponse, response.PaginationMeta, error)
	FindByID(id string) (dto.PermissionResponse, error)
	Update(id string, req dto.UpdatePermissionRequest) (dto.PermissionResponse, error)
	Delete(id string) error
}

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{repo}
}

func (s *permissionService) Create(req dto.CreatePermissionRequest) (dto.PermissionResponse, error) {
	// Cek apakah slug sudah ada
	existing, _ := s.repo.FindBySlug(req.Slug)
	if existing.ID != "" {
		return dto.PermissionResponse{}, errors.New("permission slug already exists")
	}

	permission := model.Permission{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := s.repo.Create(&permission); err != nil {
		return dto.PermissionResponse{}, err
	}

	return dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Description: permission.Description,
	}, nil
}

func (s *permissionService) FindAll() ([]dto.PermissionResponse, error) {
	permissions, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.PermissionResponse, 0)
	for _, p := range permissions {
		result = append(result, dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
		})
	}
	return result, nil
}

func (s *permissionService) FindWithFilter(filter dto.PermissionFilter) ([]dto.PermissionResponse, response.PaginationMeta, error) {
	permissions, total, err := s.repo.FindWithFilter(filter.Search, filter.Page, filter.Limit)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	result := make([]dto.PermissionResponse, 0)
	for _, p := range permissions {
		result = append(result, dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
		})
	}

	// hitung pagination meta
	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	meta := response.PaginationMeta{
		TotalRecords: total,
		TotalPages:   totalPages,
		CurrentPage:  filter.Page,
		PageSize:     filter.Limit,
	}

	return result, meta, nil
}

func (s *permissionService) FindByID(id string) (dto.PermissionResponse, error) {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return dto.PermissionResponse{}, errors.New("permission not found")
	}
	return dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Description: permission.Description,
	}, nil
}

func (s *permissionService) Update(id string, req dto.UpdatePermissionRequest) (dto.PermissionResponse, error) {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return dto.PermissionResponse{}, errors.New("permission not found")
	}

	// Cek slug uniqueness jika slug berubah
	if permission.Slug != req.Slug {
		existing, _ := s.repo.FindBySlug(req.Slug)
		if existing.ID != "" {
			return dto.PermissionResponse{}, errors.New("permission slug already exists")
		}
	}

	permission.Name = req.Name
	permission.Slug = req.Slug
	permission.Description = req.Description

	if err := s.repo.Update(&permission); err != nil {
		return dto.PermissionResponse{}, err
	}

	return dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Description: permission.Description,
	}, nil
}

func (s *permissionService) Delete(id string) error {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("permission not found")
	}
	return s.repo.Delete(&permission)
}