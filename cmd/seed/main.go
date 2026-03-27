package main

import (
	"log"

	"pos-be/config"
	"pos-be/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("failed to init db:", err)
	}

	// ensure schema
	if err := db.AutoMigrate(
		&model.Store{},
		&model.Tenant{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.User{},
	); err != nil {
		log.Fatal("migration failed:", err)
	}

	// permissions list (CRUD per resource)
	permSlugs := []string{
		"category.create", "category.read", "category.update", "category.delete",
		"product.create", "product.read", "product.update", "product.delete",
		"transaction.create", "transaction.read", "transaction.update", "transaction.delete",
		"stockmovement.create", "stockmovement.read", "stockmovement.update", "stockmovement.delete",
		"user.create", "user.read", "user.update", "user.delete",
		"role.create", "role.read", "role.update", "role.delete",
		"permission.read",
		"store.create", "store.read", "store.update", "store.delete",
		"tenant.create", "tenant.read", "tenant.update", "tenant.delete",
		"order.create", "order.read", "order.update", "order.delete",
		"payment.create", "payment.read", "payment.update", "payment.delete",
	}

	var permissions []model.Permission
	for _, slug := range permSlugs {
		p := model.Permission{Slug: slug, Name: slug}
		// use FirstOrCreate to be idempotent
		if err := db.Where("slug = ?", slug).FirstOrCreate(&p).Error; err != nil {
			log.Fatal("failed to create permission:", err)
		}
		permissions = append(permissions, p)
	}

	// create admin role
	adminRole := model.Role{Name: "Administrator", Slug: "admin", Description: "System administrator with full privileges"}
	if err := db.Where("slug = ?", adminRole.Slug).FirstOrCreate(&adminRole).Error; err != nil {
		log.Fatal("failed to create role:", err)
	}

	// link all permissions to admin role via role_permissions, idempotent
	var rolePerms []model.RolePermission
	for _, p := range permissions {
		rolePerms = append(rolePerms, model.RolePermission{RoleID: adminRole.ID, PermissionID: p.ID})
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rolePerms).Error; err != nil {
		log.Fatal("failed to link role_permissions:", err)
	}

	// create admin user
	adminEmail := "admin@local.com"
	var admin model.User
	if err := db.Where("email = ?", adminEmail).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			admin = model.User{Name: "Seed Admin", Email: adminEmail, Password: string(pass)}
			if err := db.Create(&admin).Error; err != nil {
				log.Fatal("failed to create admin user:", err)
			}
		} else {
			log.Fatal("failed to query admin user:", err)
		}
	}

	// create store with owner = admin
	store := model.Store{Name: "Seed Store", OwnerID: &admin.ID}
	if err := db.Where("name = ?", store.Name).FirstOrCreate(&store).Error; err != nil {
		log.Fatal("failed to create store:", err)
	}

	// create tenant under store
	tenant := model.Tenant{Name: "Seed Tenant", StoreID: &store.ID}
	if err := db.Where("name = ? AND store_id = ?", tenant.Name, store.ID).FirstOrCreate(&tenant).Error; err != nil {
		log.Fatal("failed to create tenant:", err)
	}

	// assign tenant to admin user (update)
	if admin.TenantID == nil || *admin.TenantID != tenant.ID {
		admin.TenantID = &tenant.ID
		if err := db.Save(&admin).Error; err != nil {
			log.Fatal("failed to update admin tenant:", err)
		}
	}

	// assign admin role to user via user_roles table (unscoped)
	ur := model.UserRole{UserID: admin.ID, RoleID: adminRole.ID}
	if err := db.Omit("tenant_id", "store_id", "created_at").Clauses(clause.OnConflict{DoNothing: true}).Create(&ur).Error; err != nil {
		log.Fatal("failed to assign role to user:", err)
	}

	log.Println("✅ Seeder finished: admin user, role, permissions, store, tenant created/ensured")
}
