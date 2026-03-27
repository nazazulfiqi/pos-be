package repository

import (
	"time"

	"pos-be/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	AssignRoles(userID string, roleIDs []string) error
	GetRoleSlugsByUser(userID string) ([]string, error)
	GetPermissionSlugsByUser(userID string) ([]string, error)
	FindByID(id string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AssignRoles(userID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return nil
	}
	var userRoles []model.UserRole
	now := time.Now()
	for _, rid := range roleIDs {
		userRoles = append(userRoles, model.UserRole{
			UserID:    userID,
			RoleID:    rid,
			CreatedAt: now,
		})
	}
	// Omit tenant_id, store_id, created_at to stay compatible with older schema
	return r.db.Omit("tenant_id", "store_id", "created_at").Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
}

func (r *userRepository) GetRoleSlugsByUser(userID string) ([]string, error) {
	var slugs []string
	rows, err := r.db.Table("roles").Select("roles.slug").Joins("join user_roles on user_roles.role_id = roles.id").Where("user_roles.user_id = ?", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

func (r *userRepository) GetPermissionSlugsByUser(userID string) ([]string, error) {
	var slugs []string
	// permissions linked via roles -> role_permissions -> permissions
	rows, err := r.db.Table("permissions").Select("distinct permissions.slug").
		Joins("join role_permissions on role_permissions.permission_id = permissions.id").
		Joins("join roles on roles.id = role_permissions.role_id").
		Joins("join user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Preload("Tenant.Store").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
