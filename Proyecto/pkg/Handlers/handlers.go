package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"proyecto/pkg/Product"
	"strconv"

	"github.com/go-chi/chi/v5"
)

/* Errors definition */
var errorOutOfIndex = errors.New("invalid ID")

/* Product handler definition */
type ProductHandler struct {
	Products []Product.TProduct
}

// NewProductHandler creates a new default valued productHandler
func NewProductHandler() ProductHandler {
	return ProductHandler{}
}

/* Internal Funtions */
// searchProduct search a product on products slice by the given id.
// searchProduct(int) -> (TProduct, error)
// Args:
//
//	id: ID of the desirable product.
//
// Return:
//
//	Product:   The product from products which matchs with the ID
//	error: 	   Error raised during the execution (if exists).
func (p *ProductHandler) searchProduct(id int) (Product.TProduct, error) {
	if id < 1 || id > len(p.Products) {
		return Product.TProduct{}, errorOutOfIndex
	}
	return p.Products[id-1], nil
}

// filterProductsByPrice filters all the products which price is below price
// filterProductsByPrice(float64) -> []Product.TProduct
func (p *ProductHandler) filterProductsByPrice(price float64) []Product.TProduct {
	var result []Product.TProduct
	for _, value := range p.Products {
		if value.Price > price {
			result = append(result, value)
		}
	}
	return result
}

// initDataSource initializes the data used in server endpoint routines
func (p *ProductHandler) InitHandler() error {
	var err error
	p.Products, err = Product.DumpJson("pkg/Product/products.json")
	if err != nil {
		return err
	}
	return nil
}

/* Endpoint function handlers */
// [GET] getAllProducts returns a slice wich contains all the products avaliable on the website
func (p *ProductHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productsJSON, err := json.Marshal(p.Products)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("An unexpected error ocurred during data loading."))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(productsJSON)
	}
}

// [GET] getProductByID search a product by ID and return if there is a match.
func (p *ProductHandler) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		product, err := p.searchProduct(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(productJson)
	}
}

// [GET] getProductByPrice returns all the products which price is higher than priceGt
func (p *ProductHandler) GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceGt, err := strconv.ParseFloat(chi.URLParam(r, "priceGt"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		filteredProducts := p.filterProductsByPrice(priceGt)
		filteredProducts_json, err := json.Marshal(filteredProducts)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(filteredProducts_json)
	}
}
