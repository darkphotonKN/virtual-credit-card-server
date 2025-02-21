package models

import (
	"time"

	"gorm.io/gorm"
)

// type for database connection values
type DBModel struct {
	DB *gorm.DB
}

// Models is the wrapper for ALL modules
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *gorm.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

// Model for Product in DB
// type Product struct {
// 	ID             int       `json:"id"`
// 	Name           string    `json:"name"`
// 	InventoryLevel int       `json:"inventory-level"`
// 	Price          int       `json:"price"`
// 	CreatedAt      time.Time `json:"-"`
// 	UpdatedAt      time.Time `json:"-"`
// }

type Product struct {
	ID             uint      `gorm: "primaryKey" json:"id"`
	Name           string    `json:"name"`
	InventoryLevel string    `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
