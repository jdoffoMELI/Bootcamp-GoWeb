package repository

import (
	"proyecto/internal"
)

type ProductMap struct {
	db     map[int]internal.TProduct // Database [product's id] -> product
	lastID int                       // Last product id used
}

// NewProductMap creates a new ProductMap instance
// NewProductMap(db map[int]internal.TProduct, lastID int) -> *ProductMap
// Args:
//		db:	 Database of products
//		lastID: Last product id used on the database
// Return:
//		*ProductMap: New ProductMap instance

func NewProductMap(db map[int]internal.TProduct, lastID int) *ProductMap {
	return &ProductMap{
		db:     db,
		lastID: lastID,
	}
}

// GetAllProducts returns the database of products
// GetAllProducts() -> []internal.TProduct
// Return:
//		internal.Tproduct: Database of products

func (p *ProductMap) GetAllProducts() []internal.TProduct {
	var productSlice []internal.TProduct
	for _, product := range p.db {
		productSlice = append(productSlice, product)
	}
	return productSlice
}

// GetProductByID returns a product by its id
// GetProductByID(id int) -> (internal.TProduct, error)
// Args:
//		id: Product id
// Return:
//		internal.TProduct: Product found in the database
//		error: 			   Error raised during the execution (if exists)

func (p *ProductMap) GetProductByID(id int) (internal.TProduct, error) {
	if product, ok := p.db[id]; !ok {
		return internal.TProduct{}, internal.ErrProductNotFound
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

func (p *ProductMap) GetProductByPriceGt(price float64) []internal.TProduct {
	var productSlice []internal.TProduct
	for _, product := range p.db {
		if product.Price > price {
			productSlice = append(productSlice, product)
		}
	}
	return productSlice
}

// productCodeExist checks if a product's code already exists in the database
// productCodeExist(code string) -> bool
// Args:
//		code: Product's code to check
// Return:
//		bool: True if the product's code already exists in the database, false otherwise

func (p *ProductMap) productCodeExist(code string) bool {
	for _, value := range p.db {
		if value.CodeValue == code {
			return true
		}
	}
	return false
}

// InsertNewProduct inserts a new product in the database
// InsertNewProduct(product internal.TProduct) -> error
// Args:
//		product: Product to insert
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductMap) InsertNewProduct(product *internal.TProduct) error {
	/* Check if the product's code already exist */
	if p.productCodeExist(product.CodeValue) {
		return internal.ErrProductCodeAlreadyExists
	}

	/* Insert the new product */
	p.lastID++                // Update the last ID used
	product.ID = p.lastID     // Update the product's ID
	p.db[p.lastID] = *product // Insert the new product
	return nil
}

// UpdateProduct updates a product it if it already exists on the repository
// UpdateProduct(product internal.TProduct) -> error
// Args:
//
//	product: Product to insert
//
// Return:
//
//	error: Error raised during the execution (if exists)
func (p *ProductMap) UpdateProduct(product *internal.TProduct) error {
	/* Check if the product exists */
	_, ok := p.db[product.ID]
	if !ok {
		return internal.ErrProductNotFound
	}

	/* Check for code value consistency */
	if p.productCodeExist(product.CodeValue) {
		return internal.ErrProductCodeAlreadyExists
	}
	/* Update the product */
	p.db[product.ID] = *product
	return nil
}

// DeleteProduct deletes a product from the database
// DeleteProduct(id int) -> error
// Args:
//		id: Product id
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductMap) DeleteProduct(id int) error {
	if _, ok := p.db[id]; !ok {
		return internal.ErrProductNotFound
	} else {
		delete(p.db, id)
		return nil
	}
}
