package product

import (
	"fmt"

	"github.com/darkphotonKN/virtual-credit-card-server/internal/models"
)

// Create Product Record
func CreateProductRecord(DBModel *models.DBModel, product models.Product) {

	result := DBModel.DB.Create(&product)

	fmt.Println("Result of creating new product record:", result)
}

// Get All Product Records
func GetProductRecords(DBModel *models.DBModel) []models.Product {
	// slice with product type for storage
	var products []models.Product
	DBModel.DB.Find(&products)

	return products
}
