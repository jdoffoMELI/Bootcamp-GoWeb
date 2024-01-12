package repository

import (
	"proyecto/internal"
	"sort"
)

type ProductMap struct {
	storage internal.ProductStorage // Storage
}

// NewProductMap creates a new ProductMap
// NewProductMap(storage internal.ProductStorage) -> *ProductMap
// Args:
//		storage: Product storage
// Return:
//		*ProductMap: New ProductMap

func NewProductMap(storage internal.ProductStorage) *ProductMap {
	return &ProductMap{storage: storage}
}

// GetAllProducts returns the database of products
// GetAllProducts() -> []internal.TProduct
// Return:
//		internal.Tproduct: Database of products

func (p *ProductMap) GetAllProducts() []internal.TProduct {
	/* Get the data from the storage */
	productMap, err := p.storage.GetAll()
	if err != nil {
		panic(err)
	}

	/* Convert the map to a slice */
	var productSlice []internal.TProduct
	for _, product := range productMap {
		productSlice = append(productSlice, product)
	}

	sort.Slice(productSlice, func(i, j int) bool {
		return productSlice[i].ID < productSlice[j].ID
	})

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
	/* Get the data from the storage */
	db, err := p.storage.GetAll()
	if err != nil {
		return internal.TProduct{}, internal.ErrStorageError
	}

	/* Check if the product exists */
	if product, ok := db[id]; !ok {
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
	/* Get the data from the storage */
	db, err := p.storage.GetAll()
	if err != nil {
		panic(err)
	}

	/* Filter the products by price */
	var productSlice []internal.TProduct
	for _, product := range db {
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
//		db:   Database of products
// Return:
//		bool: True if the product's code already exists in the database, false otherwise

func (p *ProductMap) productCodeExist(product internal.TProduct, db map[int]internal.TProduct) bool {
	for _, value := range db {
		if value.CodeValue == product.CodeValue && product.ID != value.ID {
			return true
		}
	}
	return false
}

// getNewID returns a new id for a product
// getNewID(db map[int]internal.TProduct) -> int
// Args:
//		db: Database of products
// Return:
//		int: New id for a product

func getNewID(db map[int]internal.TProduct) int {
	var lastID int
	for key := range db {
		if key > lastID {
			lastID = key
		}
	}
	return lastID + 1
}

// InsertNewProduct inserts a new product in the database
// InsertNewProduct(product internal.TProduct) -> error
// Args:
//		product: Product to insert
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductMap) InsertNewProduct(product *internal.TProduct) error {
	/* Get the data from the storage */
	db, err := p.storage.GetAll()
	if err != nil {
		return internal.ErrStorageError
	}

	/* Check if the product's code already exist */
	if p.productCodeExist(*product, db) {
		return internal.ErrProductCodeAlreadyExists
	}

	/* Insert the new product */
	newID := getNewID(db) // Get a new id for the product
	product.ID = newID    // Update the product's ID
	db[newID] = *product  // Insert the new product

	/* Save the changes in the storage */
	if err = p.storage.WriteAll(db); err != nil {
		return internal.ErrStorageError
	}

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
	/* Get the data from the storage */
	db, err := p.storage.GetAll()
	if err != nil {
		return internal.ErrStorageError
	}

	/* Check if the product exists */
	_, ok := db[product.ID]
	if !ok {
		return internal.ErrProductNotFound
	}

	/* Check for code value consistency */
	if p.productCodeExist(*product, db) {
		return internal.ErrProductCodeAlreadyExists
	}

	/* Update the product */
	db[product.ID] = *product

	/* Save the changes in the storage */
	if err = p.storage.WriteAll(db); err != nil {
		return internal.ErrStorageError
	}

	return nil
}

// DeleteProduct deletes a product from the database
// DeleteProduct(id int) -> error
// Args:
//		id: Product id
// Return:
//		error: Error raised during the execution (if exists)

func (p *ProductMap) DeleteProduct(id int) error {
	/* Get the data from the storage */
	db, err := p.storage.GetAll()
	if err != nil {
		return internal.ErrStorageError
	}

	/* Check if the product exists */
	if _, ok := db[id]; !ok {
		return internal.ErrProductNotFound
	}

	/* Delete the product */
	delete(db, id)
	if err = p.storage.WriteAll(db); err != nil {
		return internal.ErrStorageError
	}
	return nil
}
