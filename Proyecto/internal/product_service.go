package internal

import "errors"

/* Errors definition */
var (
	ErrEmptyField           = errors.New("empty field")
	ErrInvalidDate          = errors.New("invalid date")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductNotExists     = errors.New("product not exists")
)

/* Product service definition */
type ProductService interface {
	GetAllProducts() []TProduct                   // Return all the products.
	GetProductByID(id int) (TProduct, error)      // Return a product by its id.
	GetProductByPriceGt(price float64) []TProduct // Return a slice of products with a price greater than the given price.
	InsertNewProduct(product *TProduct) error     // Add a new product into the repository.
	UpdateProduct(product *TProduct) error        // Update a product from the repository if it exists.
	DeleteProduct(id int) error                   // Delete a product from the repository.
}
