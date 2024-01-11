package internal

import (
	"errors"
)

/* Errors definition */
var (
	ErrProductNotFound          = errors.New("product not found")
	ErrProductCodeAlreadyExists = errors.New("product code already exists")
)

/* Product repository definition */
type ProductRepository interface {
	GetAllProducts() []TProduct                   // Return all the products in the repository.
	GetProductByID(id int) (TProduct, error)      // Return a product by its id.
	GetProductByPriceGt(price float64) []TProduct // Return a slice of products with a price greater than the given price.
	InsertNewProduct(product *TProduct) error     // Add a new product into the repository.
}
