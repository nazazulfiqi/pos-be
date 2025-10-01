package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;size:100;not null"`
	Password  string    `gorm:"column:password_hash;size:255;not null"`
	RoleID    uint      `gorm:"column:role_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Customer struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Phone     string    `gorm:"column:phone;size:20"`
	Email     string    `gorm:"column:email;size:100"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Products  []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	Name       string    `gorm:"column:name;size:200;not null"`
	SKU        string    `gorm:"column:sku;uniqueIndex;size:50"`
	CategoryID uint      `gorm:"column:category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	Price      float64   `gorm:"column:price;not null"`
	Stock      int       `gorm:"column:stock;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type StockMovement struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	ProductID uint      `gorm:"column:product_id"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Type      string    `gorm:"column:type;size:10;not null"` // in/out
	Quantity  int       `gorm:"column:quantity;not null"`
	Note      string    `gorm:"column:note;size:255"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Transaction struct {
	ID            uint              `gorm:"primaryKey;column:id"`
	UserID        uint              `gorm:"column:user_id"`
	User          User              `gorm:"foreignKey:UserID"`
	CustomerID    *uint             `gorm:"column:customer_id"`
	Customer      *Customer         `gorm:"foreignKey:CustomerID"`
	TotalAmount   float64           `gorm:"column:total_amount;not null"`
	PaymentMethod string            `gorm:"column:payment_method;size:50;not null"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;autoUpdateTime"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            uint    `gorm:"primaryKey;column:id"`
	TransactionID uint    `gorm:"column:transaction_id;index"`
	ProductID     uint    `gorm:"column:product_id"`
	Product       Product `gorm:"foreignKey:ProductID"`
	Quantity      int     `gorm:"column:quantity;not null"`
	Price         float64 `gorm:"column:price;not null"`
	Subtotal      float64 `gorm:"column:subtotal;not null"`
}

type Order struct {
	ID          uint        `gorm:"primaryKey;column:id"`
	UserID      uint        `gorm:"column:user_id"`
	User        User        `gorm:"foreignKey:UserID"`
	CustomerID  uint        `gorm:"column:customer_id"`
	Customer    Customer    `gorm:"foreignKey:CustomerID"`
	OrderNumber string      `gorm:"column:order_number;size:50;uniqueIndex;not null"`
	Status      string      `gorm:"column:status;size:50;not null;default:'pending'"`
	TotalAmount float64     `gorm:"column:total_amount;not null"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	Payments    []Payment   `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey;column:id"`
	OrderID   uint    `gorm:"column:order_id;index"`
	ProductID uint    `gorm:"column:product_id"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"column:quantity;not null"`
	Price     float64 `gorm:"column:price;not null"`
	Subtotal  float64 `gorm:"column:subtotal;not null"`
}

type Payment struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	OrderID               uint      `gorm:"column:order_id;index"`
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

type WebhookLog struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	Provider   string    `gorm:"column:provider;size:50;not null"`
	Payload    string    `gorm:"column:payload;type:text"`
	Headers    string    `gorm:"column:headers;type:text"`
	ReceivedAt time.Time `gorm:"column:received_at;autoCreateTime"`
}
