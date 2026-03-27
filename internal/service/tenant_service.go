package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"pos-be/internal/response"
)

type TenantService interface {
	Create(req dto.CreateTenantRequest) (dto.TenantResponse, error)
	FindAll(storeID *string) ([]dto.TenantResponse, error)
	FindWithFilter(filter dto.TenantFilter) ([]dto.TenantResponse, response.PaginationMeta, error)
	FindByID(id string) (dto.TenantResponse, error)
	Update(id string, req dto.UpdateTenantRequest) (dto.TenantResponse, error)
	Delete(id string) error
}

type tenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) TenantService {
	return &tenantService{repo}
}

func (s *tenantService) Create(req dto.CreateTenantRequest) (dto.TenantResponse, error) {
	t := model.Tenant{
		Name:    req.Name,
		StoreID: req.StoreID,
	}
	if err := s.repo.Create(&t); err != nil {
		return dto.TenantResponse{}, err
	}
	return dto.TenantResponse{ID: t.ID, Name: t.Name, StoreID: t.StoreID}, nil
}

func (s *tenantService) FindAll(storeID *string) ([]dto.TenantResponse, error) {
	tenants, err := s.repo.FindAll(storeID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.TenantResponse, 0)
	for _, t := range tenants {
		res = append(res, dto.TenantResponse{ID: t.ID, Name: t.Name, StoreID: t.StoreID})
	}
	return res, nil
}

func (s *tenantService) FindByID(id string) (dto.TenantResponse, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return dto.TenantResponse{}, errors.New("tenant not found")
	}
	return dto.TenantResponse{ID: t.ID, Name: t.Name, StoreID: t.StoreID}, nil
}

func (s *tenantService) Update(id string, req dto.UpdateTenantRequest) (dto.TenantResponse, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return dto.TenantResponse{}, errors.New("tenant not found")
	}
	t.Name = req.Name
	t.StoreID = req.StoreID
	if err := s.repo.Update(&t); err != nil {
		return dto.TenantResponse{}, err
	}
	return dto.TenantResponse{ID: t.ID, Name: t.Name, StoreID: t.StoreID}, nil
}

func (s *tenantService) FindWithFilter(filter dto.TenantFilter) ([]dto.TenantResponse, response.PaginationMeta, error) {
	tenants, total, err := s.repo.FindWithFilter(filter.Name, filter.StoreID, filter.Page, filter.Limit)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]dto.TenantResponse, 0)
	for _, t := range tenants {
		res = append(res, dto.TenantResponse{ID: t.ID, Name: t.Name, StoreID: t.StoreID})
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	meta := response.PaginationMeta{
		TotalRecords: total,
		TotalPages:   totalPages,
		CurrentPage:  filter.Page,
		PageSize:     filter.Limit,
	}

	return res, meta, nil
}

func (s *tenantService) Delete(id string) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("tenant not found")
	}
	return s.repo.Delete(&t)
}
