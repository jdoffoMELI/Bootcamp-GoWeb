package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	Product "proyecto/internal/product"
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
// searchProduct(id int) -> (TProduct, error)
// Args:
//		id: 	   ID of the desirable product.
// Return:
//		Product:   The product from products which matchs with the ID
//		error: 	   Error raised during the execution (if exists).

func (p *ProductHandler) searchProduct(id int) (Product.TProduct, error) {
	if id < 1 || id > len(p.Products) {
		return Product.TProduct{}, errorOutOfIndex
	}
	return p.Products[id-1], nil
}

// filterProductsByPrice filters all the products which price is below price
// filterProductsByPrice(price float64) -> []Product.TProduct
// Args:
//		price float64: 		The price to filter by.
// Return:
//		[]Product.TProduct: Slice of products which price is over "price".

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
// initDataSource() -> error
func (p *ProductHandler) InitHandler() error {
	var err error
	p.Products, err = Product.DumpJson("docs/db/products.json")
	if err != nil {
		return err
	}
	return nil
}

// existProductCode checks if there is a product with the given code value
// existProductCode(string) -> bool
func (p *ProductHandler) existProductCode(code string) bool {
	for _, value := range p.Products {
		if value.CodeValue == code {
			return true
		}
	}
	return false
}

/* Endpoint function handlers */
// [GET] getAllProducts returns a slice wich contains all the products avaliable on the website
// Url params: none

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
// URL params:
//		id (Numeric): ID of the desirable product.

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
// URL params:
//		priceGt (Numeric): Price to filter by.

func (p *ProductHandler) GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceGt, err := strconv.ParseFloat(r.URL.Query().Get("priceGt"), 64)
		// priceGt, err := strconv.ParseFloat(chi.URLParam(r, "priceGt"), 64)
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

// [POST] postNewProduct creates a new product on the website
// URL params : none
// Body params:
//	  	 name (string): 		Product name.
//		 quantity (numeric): 	Product quantity.
//		 code_value (string):	Product code. Its must be unique.
//		 expiration (string):	Product expiration date. Its must be in the format "DD/MM/YYYY".
//		 price (numeric):		Product price.
//		 is_published (bool):	(Optional) Product publication status. Default value is false.

func (p *ProductHandler) PostNewProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product.TProduct
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		/* Checks for empty values except is_published */
		if product.HasEmptyValues() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Empty values are not allowed."))
			return
		}

		/* Checks for code_value uniqueness */
		if p.existProductCode(product.CodeValue) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Product code already exists."))
			return
		}

		/* Checks for date correctness */
		if !product.HasValidDate() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid date."))
			return
		}
		product.ID = len(p.Products) + 1
		p.Products = append(p.Products, product)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Product created successfully.",
			"data":    product,
		})

	}
}
