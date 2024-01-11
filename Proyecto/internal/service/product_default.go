package service

import (
	"fmt"
	"proyecto/internal"
	"strconv"
	"strings"
)

type ProductServiceDefault struct {
	repository internal.ProductRepository
}

// NewProductServiceDefault creates a new ProductServiceDefault instance
// NewProductServiceDefault(rp internal.ProductRepository) -> *ProductServiceDefault
// Args:
//		rp: Product Repository where the service will get the data from
// Return:
//		*ProductServiceDefault: New ProductServiceDefault instance

func NewProductServiceDefault(rp internal.ProductRepository) *ProductServiceDefault {
	return &ProductServiceDefault{
		repository: rp,
	}
}

// GetAllProducts returns all the products in the repository
// GetAllProducts() -> []internal.TProduct
// Return:
//		[]internal.TProduct: Slice of products

func (p *ProductServiceDefault) GetAllProducts() []internal.TProduct {
	return p.repository.GetAllProducts()
}

// GetProductByID returns a product by its id
// GetProductByID(id int) -> (internal.TProduct, error)
// Args:
//		id: Product id
// Return:
//		internal.TProduct: Product found in the repository
//		error: 			   Error raised during the execution (if exists)

func (p *ProductServiceDefault) GetProductByID(id int) (internal.TProduct, error) {
	/* Get the product by its ID */
	if product, err := p.repository.GetProductByID(id); err == internal.ErrProductNotFound {
		return internal.TProduct{}, internal.ErrProductNotExists
	} else if err != nil {
		return internal.TProduct{}, err
	} else {
		return product, nil
	}
}

// GetProductByPriceGt returns a slice of products with a price greater than the given price
// GetProductByPriceGt(price float64) -> []internal.TProduct
// Args:
//		price: Price to compare
// Return:
//		[]internal.TProduct: Slice of products with a price greater than the given price

func (p *ProductServiceDefault) GetProductByPriceGt(price float64) []internal.TProduct {
	return p.repository.GetProductByPriceGt(price)
}

// EmptyValues checks if the product has empty values
// EmptyValues(p internal.TProduct) -> []string
// Args:
//		p: Product to check
// Return:
//		[]string: Slice of empty value field names

func EmptyValues(p internal.TProduct) []string {
	var emptyFields = make([]string, 0)
	if p.Name == "" {
		emptyFields = append(emptyFields, "Name")
	}

	if p.Quantity == 0 {
		emptyFields = append(emptyFields, "Quantity")
	}

	if p.CodeValue == "" {
		emptyFields = append(emptyFields, "Code Value")
	}

	if p.Expiration == "" {
		emptyFields = append(emptyFields, "Expiration")
	}

	if p.Price == 0.0 {
		emptyFields = append(emptyFields, "Price")
	}

	return emptyFields
}

// validateDate checks if the date is valid. It must have the format dd/mm/yyyy and the date must be valid
// validateDate(date string) -> bool
// Args:
//		date: Date to validate
// Return:
//		bool: True if the date is valid, false otherwise

func validateDate(date string) bool {
	/* Decontruct the date */
	tokenSlice := strings.Split(date, "/")
	if len(tokenSlice) != 3 {
		return false
	}
	tokenDay, tokenMonth, tokenYear := tokenSlice[0], tokenSlice[1], tokenSlice[2]

	/* Parse its components into numeric values */
	day, err := strconv.Atoi(tokenDay)
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(tokenMonth)
	if err != nil {
		return false
	}
	year, err := strconv.Atoi(tokenYear)
	if err != nil {
		return false
	}

	/* Validate the date */
	return day > 0 && day <= 31 && month > 0 && month <= 12 && year > 1900 && year <= 2024
}

// InsertNewProduct inserts a new product into the repository
// InsertNewProduct(product internal.TProduct) -> error
// Args:
//		product: Product to insert
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductServiceDefault) InsertNewProduct(product *internal.TProduct) error {
	/* Empty fields validation */
	if emptyFields := EmptyValues(*product); len(emptyFields) != 0 {
		return fmt.Errorf("%w: %s", internal.ErrEmptyField, strings.Join(emptyFields, ", "))
	}
	/* Date validation */
	if !validateDate(product.Expiration) {
		return internal.ErrInvalidDate
	}

	/* Insert the new product into the repository */
	if err := p.repository.InsertNewProduct(product); err == internal.ErrProductCodeAlreadyExists {
		return internal.ErrProductAlreadyExists
	} else {
		return err
	}
}

// UpdateProduct updates a product it if it already exists
// UpdateProduct(product internal.TProduct) -> error
// Args:
//		product: Product to insert or update
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductServiceDefault) UpdateProduct(product *internal.TProduct) error {
	/* Empty fields validation */
	if emptyFields := EmptyValues(*product); len(emptyFields) != 0 {
		return fmt.Errorf("%w: %s", internal.ErrEmptyField, strings.Join(emptyFields, ", "))
	}

	/* Date validation */
	if !validateDate(product.Expiration) {
		return internal.ErrInvalidDate
	}

	/* Insert the new product into the repository */
	if err := p.repository.UpdateProduct(product); err == internal.ErrProductCodeAlreadyExists {
		return internal.ErrProductAlreadyExists
	} else {
		return err
	}
}

// DeleteProduct deletes a product from the repository
// DeleteProduct(id int) -> error
// Args:
//		id: Product id
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductServiceDefault) DeleteProduct(id int) error {
	/* Delete the product from the repository */
	if err := p.repository.DeleteProduct(id); err == internal.ErrProductNotFound {
		return internal.ErrProductNotExists
	} else {
		return err
	}
}
