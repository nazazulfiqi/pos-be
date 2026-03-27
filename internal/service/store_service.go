package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"pos-be/internal/response"
)

type StoreService interface {
	Create(req dto.CreateStoreRequest) (dto.StoreResponse, error)
	FindAll(ownerID *string) ([]dto.StoreResponse, error)
	FindWithFilter(filter dto.StoreFilter) ([]dto.StoreResponse, response.PaginationMeta, error)
	FindByID(id string) (dto.StoreResponse, error)
	Update(id string, req dto.UpdateStoreRequest) (dto.StoreResponse, error)
	Delete(id string) error
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo}
}

func (s *storeService) Create(req dto.CreateStoreRequest) (dto.StoreResponse, error) {
	store := model.Store{
		Name:    req.Name,
		OwnerID: req.OwnerID,
	}
	if err := s.repo.Create(&store); err != nil {
		return dto.StoreResponse{}, err
	}
	return dto.StoreResponse{ID: store.ID, Name: store.Name, OwnerID: store.OwnerID}, nil
}

func (s *storeService) FindAll(ownerID *string) ([]dto.StoreResponse, error) {
	stores, err := s.repo.FindAll(ownerID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.StoreResponse, 0)
	for _, st := range stores {
		res = append(res, dto.StoreResponse{ID: st.ID, Name: st.Name, OwnerID: st.OwnerID})
	}
	return res, nil
}

func (s *storeService) FindWithFilter(filter dto.StoreFilter) ([]dto.StoreResponse, response.PaginationMeta, error) {
	stores, total, err := s.repo.FindWithFilter(filter.Name, filter.OwnerID, filter.Page, filter.Limit)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]dto.StoreResponse, 0)
	for _, st := range stores {
		res = append(res, dto.StoreResponse{ID: st.ID, Name: st.Name, OwnerID: st.OwnerID})
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

func (s *storeService) FindByID(id string) (dto.StoreResponse, error) {
	st, err := s.repo.FindByID(id)
	if err != nil {
		return dto.StoreResponse{}, errors.New("store not found")
	}
	return dto.StoreResponse{ID: st.ID, Name: st.Name, OwnerID: st.OwnerID}, nil
}

func (s *storeService) Update(id string, req dto.UpdateStoreRequest) (dto.StoreResponse, error) {
	st, err := s.repo.FindByID(id)
	if err != nil {
		return dto.StoreResponse{}, errors.New("store not found")
	}
	st.Name = req.Name
	st.OwnerID = req.OwnerID
	if err := s.repo.Update(&st); err != nil {
		return dto.StoreResponse{}, err
	}
	return dto.StoreResponse{ID: st.ID, Name: st.Name, OwnerID: st.OwnerID}, nil
}

func (s *storeService) Delete(id string) error {
	st, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("store not found")
	}
	return s.repo.Delete(&st)
}
