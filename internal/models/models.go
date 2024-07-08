package models

import (
	"context"
	"database/sql"
	"time"
)

// type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for ALL modules
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

// Model for Product in DB
type Product struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	InventoryLevel int       `json:"inventory-level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Add a method for querying the DB for Products
func (m *DBModel) GetProduct(id string) (Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // cancel when the function returns

	// query db
	var productQueried Product
	row := m.DB.QueryRowContext(ctx, "select id, name from widgets where id = ?", id)
	err := row.Scan(&productQueried.ID, &productQueried.Name)

	if err != nil {
		return productQueried, err
	}

	return productQueried, nil
}
