package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"
	"pos-be/internal/response"
)

type RoleService interface {
	Create(req dto.CreateRoleRequest) (dto.RoleResponse, error)
	FindAll(tenantID *string) ([]dto.RoleResponse, error)
	FindWithFilter(filter dto.RoleFilter) ([]dto.RoleResponse, response.PaginationMeta, error)
	FindByID(id string) (dto.RoleResponse, error)
	Update(id string, req dto.UpdateRoleRequest) (dto.RoleResponse, error)
	Delete(id string) error
}

type roleService struct {
	repo           repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(repo repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{repo, permissionRepo}
}

func (s *roleService) buildRoleResponse(role model.Role) dto.RoleResponse {
	permissions := make([]dto.PermissionResponse, 0)
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
		})
	}

	return dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		TenantID:    role.TenantID,
		StoreID:     role.StoreID,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func (s *roleService) Create(req dto.CreateRoleRequest) (dto.RoleResponse, error) {
	// Cek apakah slug sudah ada
	existing, _ := s.repo.FindBySlug(req.Slug)
	if existing.ID != "" {
		return dto.RoleResponse{}, errors.New("role slug already exists")
	}

	// Ambil permissions berdasarkan IDs
	var permissions []model.Permission
	if len(req.PermissionIDs) > 0 {
		for _, pid := range req.PermissionIDs {
			p, err := s.permissionRepo.FindByID(pid)
			if err != nil {
				return dto.RoleResponse{}, errors.New("permission not found: " + pid)
			}
			permissions = append(permissions, p)
		}
	}

	role := model.Role{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		TenantID:    req.TenantID,
		StoreID:     req.StoreID,
		Permissions: permissions,
	}

	if err := s.repo.Create(&role); err != nil {
		return dto.RoleResponse{}, err
	}

	// Reload untuk mendapatkan data lengkap termasuk permissions
	role, _ = s.repo.FindByID(role.ID)

	return s.buildRoleResponse(role), nil
}

func (s *roleService) FindAll(tenantID *string) ([]dto.RoleResponse, error) {
	roles, err := s.repo.FindAll(tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.RoleResponse, 0)
	for _, r := range roles {
		result = append(result, s.buildRoleResponse(r))
	}
	return result, nil
}

func (s *roleService) FindWithFilter(filter dto.RoleFilter) ([]dto.RoleResponse, response.PaginationMeta, error) {
	roles, total, err := s.repo.FindWithFilter(filter.Search, filter.Page, filter.Limit, filter.TenantID, filter.StoreID)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	result := make([]dto.RoleResponse, 0)
	for _, r := range roles {
		result = append(result, s.buildRoleResponse(r))
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

func (s *roleService) FindByID(id string) (dto.RoleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return dto.RoleResponse{}, errors.New("role not found")
	}
	return s.buildRoleResponse(role), nil
}

func (s *roleService) Update(id string, req dto.UpdateRoleRequest) (dto.RoleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return dto.RoleResponse{}, errors.New("role not found")
	}

	// Cek slug uniqueness jika slug berubah
	if role.Slug != req.Slug {
		existing, _ := s.repo.FindBySlug(req.Slug)
		if existing.ID != "" {
			return dto.RoleResponse{}, errors.New("role slug already exists")
		}
	}

	// Ambil permissions berdasarkan IDs
	var permissions []model.Permission
	if len(req.PermissionIDs) > 0 {
		for _, pid := range req.PermissionIDs {
			p, err := s.permissionRepo.FindByID(pid)
			if err != nil {
				return dto.RoleResponse{}, errors.New("permission not found: " + pid)
			}
			permissions = append(permissions, p)
		}
	}

	role.Name = req.Name
	role.Slug = req.Slug
	role.Description = req.Description
	role.Permissions = permissions

	if err := s.repo.Update(&role); err != nil {
		return dto.RoleResponse{}, err
	}

	// Reload untuk mendapatkan data lengkap
	role, _ = s.repo.FindByID(role.ID)

	return s.buildRoleResponse(role), nil
}

func (s *roleService) Delete(id string) error {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}
	return s.repo.Delete(&role)
}