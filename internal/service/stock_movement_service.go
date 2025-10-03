package service

import (
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"time"
)

type StockMovementService interface {
	Create(req dto.StockMovementCreateRequest) (dto.StockMovementResponse, error)
	FindAll() ([]dto.StockMovementResponse, error)
	FindByProduct(productID uint) ([]dto.StockMovementResponse, error)
}

type stockMovementService struct {
	repo repository.StockMovementRepository
}

func NewStockMovementService(repo repository.StockMovementRepository) StockMovementService {
	return &stockMovementService{repo}
}

func (s *stockMovementService) Create(req dto.StockMovementCreateRequest) (dto.StockMovementResponse, error) {
	movement := model.StockMovement{
		ProductID:     req.ProductID,
		Type:          req.Type,
		Quantity:      req.Quantity,
		Note:          req.Note,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
		CreatedAt:     time.Now(),
	}

	err := s.repo.Create(&movement)
	if err != nil {
		return dto.StockMovementResponse{}, err
	}

	return dto.StockMovementResponse{
		ID:            movement.ID,
		ProductID:     movement.ProductID,
		ProductName:   movement.Product.Name,
		Type:          movement.Type,
		Quantity:      movement.Quantity,
		Note:          movement.Note,
		ReferenceID:   movement.ReferenceID,
		ReferenceType: movement.ReferenceType,
		CreatedAt:     movement.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *stockMovementService) FindAll() ([]dto.StockMovementResponse, error) {
	movements, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.StockMovementResponse, 0)
	for _, m := range movements {
		result = append(result, dto.StockMovementResponse{
			ID:            m.ID,
			ProductID:     m.ProductID,
			ProductName:   m.Product.Name,
			Type:          m.Type,
			Quantity:      m.Quantity,
			Note:          m.Note,
			ReferenceID:   m.ReferenceID,
			ReferenceType: m.ReferenceType,
			CreatedAt:     m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (s *stockMovementService) FindByProduct(productID uint) ([]dto.StockMovementResponse, error) {
	movements, err := s.repo.FindByProduct(productID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.StockMovementResponse, 0)
	for _, m := range movements {
		result = append(result, dto.StockMovementResponse{
			ID:            m.ID,
			ProductID:     m.ProductID,
			ProductName:   m.Product.Name,
			Type:          m.Type,
			Quantity:      m.Quantity,
			Note:          m.Note,
			ReferenceID:   m.ReferenceID,
			ReferenceType: m.ReferenceType,
			CreatedAt:     m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}
