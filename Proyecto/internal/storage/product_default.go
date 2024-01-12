package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"proyecto/internal"
)

// ProductStorageDefault is the default implementation of ProductStorage
type ProductStorageDefault struct {
	filePath string // File path
}

// NewProductStorageDefault creates a new ProductStorageDefault
// NewProductStorageDefault(filePath string) -> *ProductStorageDefault
// Args:
// 	filePath string: File path
// Returns:
// 	*ProductStorageDefault: New ProductStorageDefault

func NewProductStorageDefault(filePath string) *ProductStorageDefault {
	return &ProductStorageDefault{filePath: filePath}
}

// DumpJson creates a slice of products from a json file
// DumpJson(string) -> ([]TProduct, error)
// Args:
//		jsonPath: Json file path.
// Return:
//		[]Product: Slice of products retrieved from a json file.
//		error: 	   Error raised during the execution (if exists).

func dumpJson(jsonPath string) ([]internal.TProduct, error) {
	var jsonSlice []internal.TProduct
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	for {
		if err := jsonDecoder.Decode(&jsonSlice); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return jsonSlice, nil
}

// SliceToMap creates a map of products from a slice of products
// SliceToMap([]TProduct) -> map[int]TProduct
// Args:
//		slice: Slice of products.
// Return:
//		map[int]TProduct: Map of products.

func sliceToMap(slice []internal.TProduct) map[int]internal.TProduct {
	m := make(map[int]internal.TProduct)
	for _, v := range slice {
		m[v.ID] = v
	}
	return m
}

// GetAll gets all the products from the storage
// GetAll() -> (map[int]TProduct, error)
// Return:
//		map[int]TProduct: Map of products.
//		error: 		   Error raised during the execution (if exists).

func (p *ProductStorageDefault) GetAll() (map[int]internal.TProduct, error) {
	/* Dumpt all the products into memory */
	data, err := dumpJson(p.filePath)
	if err != nil {
		return nil, internal.ErrBadFile
	}
	/* Convert []TProduct -> map[int]TProduct */
	return sliceToMap(data), nil
}

// WriteAll writes all the products to the storage
// WriteAll(map[int]TProduct) -> error
// Args:
//		products: Map of products.
// Return:
//		error: Error raised during the execution (if exists).

func (p *ProductStorageDefault) WriteAll(products map[int]internal.TProduct) error {
	/* Open a file descriptor */
	file, err := os.Create(p.filePath)
	if err != nil {
		return internal.ErrBadFile
	}
	defer file.Close()

	/* Write all the products into the storage */
	var productsByte []internal.TProduct
	for _, value := range products {
		productsByte = append(productsByte, value)
	}
	return json.NewEncoder(file).Encode(productsByte)
}
