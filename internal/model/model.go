package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36)"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;size:100;not null"`
	Password  string    `gorm:"column:password_hash;size:255;not null"`
	TenantID  *string   `gorm:"column:tenant_id"`
	Tenant    *Tenant   `gorm:"foreignKey:TenantID"`
	Roles     []Role    `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

// Store groups multiple tenants (e.g., a brand or account owner)
type Store struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36)"`
	Name      string    `gorm:"column:name;size:200;not null"`
	OwnerID   *string   `gorm:"column:owner_id"` // optional reference to a User who owns the store
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Tenants   []Tenant  `gorm:"foreignKey:StoreID"`
}

func (s *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

// Tenant represents a customer/organization (e.g., an outlet/cabang) using the SaaS application.
type Tenant struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36)"`
	StoreID   *string   `gorm:"column:store_id"`
	Store     *Store    `gorm:"foreignKey:StoreID"`
	Name      string    `gorm:"column:name;size:200;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Users     []User    `gorm:"foreignKey:TenantID"`
	Roles     []Role    `gorm:"foreignKey:TenantID"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}

// Role groups permissions. Roles can be tenant-scoped (TenantID != nil) or global (TenantID == nil).
type Role struct {
	ID          string       `gorm:"primaryKey;column:id;type:char(36)"`
	Name        string       `gorm:"column:name;size:100;not null"`
	Slug        string       `gorm:"column:slug;size:100;not null;index"`
	Description string       `gorm:"column:description;size:255"`
	TenantID    *string      `gorm:"column:tenant_id"`
	StoreID     *string      `gorm:"column:store_id"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}

// Permission represents a single action/ability in the system.
type Permission struct {
	ID          string    `gorm:"primaryKey;column:id;type:char(36)"`
	Name        string    `gorm:"column:name;size:150;not null"`
	Slug        string    `gorm:"column:slug;size:150;not null;index"`
	Description string    `gorm:"column:description;size:255"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Roles       []Role    `gorm:"many2many:role_permissions;"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

// UserRole is an optional explicit join model for user_roles.
type UserRole struct {
	UserID    string    `gorm:"column:user_id;primaryKey;type:char(36)"`
	RoleID    string    `gorm:"column:role_id;primaryKey;type:char(36)"`
	TenantID  *string   `gorm:"column:tenant_id"`
	StoreID   *string   `gorm:"column:store_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// RolePermission is an optional explicit join model for role_permissions.
type RolePermission struct {
	RoleID       string    `gorm:"column:role_id;primaryKey;type:char(36)"`
	PermissionID string    `gorm:"column:permission_id;primaryKey;type:char(36)"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Customer struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID  *string   `gorm:"column:tenant_id"`
	Tenant    *Tenant   `gorm:"foreignKey:TenantID"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Phone     string    `gorm:"column:phone;size:20"`
	Email     string    `gorm:"column:email;size:100"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type Category struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID  *string   `gorm:"column:tenant_id"`
	Tenant    *Tenant   `gorm:"foreignKey:TenantID"`
	Name      string    `gorm:"column:name;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Products  []Product `gorm:"foreignKey:CategoryID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type Product struct {
	ID         string    `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID   *string   `gorm:"column:tenant_id"`
	Tenant     *Tenant   `gorm:"foreignKey:TenantID"`
	Name       string    `gorm:"column:name;size:200;not null"`
	SKU        string    `gorm:"column:sku;uniqueIndex;size:50"`
	CategoryID string    `gorm:"column:category_id;type:char(36)"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	Price      float64   `gorm:"column:price;not null"`
	Stock      int       `gorm:"column:stock;not null"`
	ImageURL   string    `gorm:"column:image_url;size:255"`
	PublicID   string    `gorm:"column:image_public_id;size:150"` // NEW
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

type StockMovement struct {
	ID        string  `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID  *string `gorm:"column:tenant_id"`
	Tenant    *Tenant `gorm:"foreignKey:TenantID"`
	ProductID string  `gorm:"column:product_id;type:char(36)"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Type      string  `gorm:"column:type;size:10;not null"` // in/out
	Quantity  int     `gorm:"column:quantity;not null"`
	Note      string  `gorm:"column:note;size:255"`

	// Reference polymorphic
	ReferenceID   string `gorm:"column:reference_id;size:50"`   // TRX-25-0001, ORD-25-0001, dll
	ReferenceType string `gorm:"column:reference_type;size:50"` // transaction, order, adjustment, etc

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (s *StockMovement) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

// Payment status constants
const (
	PaymentStatusPending uint = 0
	PaymentStatusSuccess uint = 1
	PaymentStatusCancel  uint = 2
)

// Payment method constants
const (
	PaymentMethodCash     uint = 1 // Cash
	PaymentMethodMidtrans uint = 2 // Midtrans
)

type Transaction struct {
	ID            string            `gorm:"primaryKey;column:id;size:20"`
	UserID        string            `gorm:"column:user_id;type:char(36)"`
	User          User              `gorm:"foreignKey:UserID"`
	TenantID      string            `gorm:"column:tenant_id;not null"`
	Tenant        *Tenant           `gorm:"foreignKey:TenantID"`
	CustomerID    *string           `gorm:"column:customer_id"`
	Customer      *Customer         `gorm:"foreignKey:CustomerID"`
	TotalAmount   float64           `gorm:"column:total_amount;not null"`
	PaymentMethod uint              `gorm:"column:payment_method"` // 1=Cash, 2=Midtrans
	PaymentStatus uint              `gorm:"column:payment_status;not null;default:0"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;autoUpdateTime"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            string  `gorm:"primaryKey;column:id;type:char(36)"`
	TransactionID string  `gorm:"column:transaction_id;index"`
	TenantID      string  `gorm:"column:tenant_id;not null"`
	ProductID     string  `gorm:"column:product_id;type:char(36)"`
	Product       Product `gorm:"foreignKey:ProductID"`
	Quantity      int     `gorm:"column:quantity;not null"`
	Price         float64 `gorm:"column:price;not null"`
	Subtotal      float64 `gorm:"column:subtotal;not null"`
}

func (ti *TransactionItem) BeforeCreate(tx *gorm.DB) (err error) {
	if ti.ID == "" {
		ti.ID = uuid.New().String()
	}
	return
}

type Order struct {
	ID          string      `gorm:"primaryKey;column:id;type:char(36)"`
	UserID      string      `gorm:"column:user_id;type:char(36)"`
	User        User        `gorm:"foreignKey:UserID"`
	TenantID    *string     `gorm:"column:tenant_id"`
	Tenant      *Tenant     `gorm:"foreignKey:TenantID"`
	CustomerID  string      `gorm:"column:customer_id;type:char(36)"`
	Customer    Customer    `gorm:"foreignKey:CustomerID"`
	OrderNumber string      `gorm:"column:order_number;size:50;uniqueIndex;not null"`
	Status      string      `gorm:"column:status;size:50;not null;default:'pending'"`
	TotalAmount float64     `gorm:"column:total_amount;not null"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	Payments    []Payment   `gorm:"foreignKey:OrderID"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return
}

type OrderItem struct {
	ID        string  `gorm:"primaryKey;column:id;type:char(36)"`
	OrderID   string  `gorm:"column:order_id;index;type:char(36)"`
	TenantID  *string `gorm:"column:tenant_id"`
	ProductID string  `gorm:"column:product_id;type:char(36)"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"column:quantity;not null"`
	Price     float64 `gorm:"column:price;not null"`
	Subtotal  float64 `gorm:"column:subtotal;not null"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == "" {
		oi.ID = uuid.New().String()
	}
	return
}

type Payment struct {
	ID                    string    `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID              *string   `gorm:"column:tenant_id"`
	Tenant                *Tenant   `gorm:"foreignKey:TenantID"`
	OrderID               string    `gorm:"column:order_id;index;type:char(36)"`
	Order                 Order     `gorm:"foreignKey:OrderID"`
	Provider              string    `gorm:"column:provider;size:50;not null"`
	ProviderTransactionID string    `gorm:"column:provider_transaction_id;size:100"`
	Method                string    `gorm:"column:method;size:50"`
	Amount                float64   `gorm:"column:amount;not null"`
	Currency              string    `gorm:"column:currency;size:10;default:'IDR'"`
	Status                string    `gorm:"column:status;size:50;not null"`
	RawResponse           string    `gorm:"column:raw_response;type:text"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

type WebhookLog struct {
	ID         string    `gorm:"primaryKey;column:id;type:char(36)"`
	TenantID   *string   `gorm:"column:tenant_id"`
	Tenant     *Tenant   `gorm:"foreignKey:TenantID"`
	Provider   string    `gorm:"column:provider;size:50;not null"`
	Payload    string    `gorm:"column:payload;type:text"`
	Headers    string    `gorm:"column:headers;type:text"`
	ReceivedAt time.Time `gorm:"column:received_at;autoCreateTime"`
}

func (w *WebhookLog) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return
}
