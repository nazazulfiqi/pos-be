package service

import (
	"errors"
	"pos-be/internal/dto"
	"pos-be/internal/model"
	"pos-be/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (dto.UserResponse, error)
	GetMe(userID string) (dto.MeResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (dto.UserResponse, error) {
	// cek apakah email sudah terdaftar
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return dto.UserResponse{}, errors.New("email already exists")
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		TenantID: req.TenantID,
	}

	if err := s.repo.Create(&user); err != nil {
		return dto.UserResponse{}, err
	}

	// assign roles
	if len(req.RoleIDs) > 0 {
		if err := s.repo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return dto.UserResponse{}, err
		}
	}

	return dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Roles:    req.RoleIDs,
		TenantID: user.TenantID,
	}, nil
}

func (s *userService) GetMe(userID string) (dto.MeResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return dto.MeResponse{}, err
	}

	roles, err := s.repo.GetRoleSlugsByUser(userID)
	if err != nil {
		return dto.MeResponse{}, err
	}
	perms, err := s.repo.GetPermissionSlugsByUser(userID)
	if err != nil {
		return dto.MeResponse{}, err
	}

	var tenant dto.TenantDTO
	var store dto.StoreDTO
	if user.Tenant != nil {
		tenant = dto.TenantDTO{ID: user.Tenant.ID, Name: user.Tenant.Name, StoreID: user.Tenant.StoreID}
		if user.Tenant.StoreID != nil && user.Tenant.Store != nil {
			store = dto.StoreDTO{ID: user.Tenant.Store.ID, Name: user.Tenant.Store.Name, OwnerID: user.Tenant.Store.OwnerID}
		}
	}

	return dto.MeResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Roles:       roles,
		Permissions: perms,
		Tenant:      &tenant,
		Store:       &store,
	}, nil
}
