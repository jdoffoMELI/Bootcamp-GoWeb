package handlers

import (
	"errors"
	"net/http"
	"os"
	"proyecto/internal"
	"proyecto/platform/web/request"
	"proyecto/platform/web/response"
	"strconv"

	"github.com/go-chi/chi/v5"
)

/* Product handler definition */
type ProductHandler struct {
	ProductService internal.ProductService // Product service instance
}

/* product field names from JSON format */
var productFields = []string{"id", "name", "quantity", "code_value", "is_published", "expiration", "price"}

// NewProductHandler creates a new default valued productHandler
// NewProductHandler(ps internal.ProductService) -> *ProductHandler
// Args:
//		ps: Product service instance
// Return:
//		*ProductHandler: New ProductHandler instance

func NewProductHandler(ps internal.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: ps,
	}
}

// ProductJSON is the JSON representation of a product
type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

// BodyRequestProductJSON is the body request for a product in JSON format
type BodyRequestProductJSON struct {
	Name        string  `json:"name"`         // Product name.
	Quantity    int     `json:"quantity"`     // Product quantity.
	CodeValue   string  `json:"code_value"`   // Product code value.
	IsPublished bool    `json:"is_published"` // Product is published (Optional)
	Expiration  string  `json:"expiration"`   // Product expiration date. Format DD/MM/YYYY
	Price       float64 `json:"price"`        // Product price.
}

// isAuthenthicated returns true if the user is authenthicated using a token
// isAuthenthicated(*http.Request) -> bool
// Args:
//		r: HTTP request
// Return:
//		bool: True if the user is authenthicated, false otherwise

func (p *ProductHandler) isAuthenthicated(r *http.Request) bool {
	userToken := r.Header.Get("TOKEN")
	return userToken == os.Getenv("TOKEN")
}

/* Endpoint function handlers */

// GetAllProducts returns a slice wich contains all the products avaliable on the website
// Url params: none
func (p *ProductHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productsJSON := p.ProductService.GetAllProducts()
		/* Send to the client all the products */
		response.JSON(w, http.StatusOK, productsJSON)
	}
}

// GetProductByID search a product by ID and return if there is a match.
// URL params:
//
//	id (Numeric): ID of the desirable product.
func (p *ProductHandler) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Retrieve the id from the url */
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid ID.")
			return
		}

		/* Search the product by id */
		product, err := p.ProductService.GetProductByID(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotExists):
				response.Text(w, http.StatusNotFound, "Product not found.")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}

		/* Send the product as response */
		response.JSON(w, http.StatusOK, product)
	}
}

// GetProductByPrice returns all the products which price is higher than priceGt
// URL params:
//
//	priceGt (Numeric): Price to filter by.
func (p *ProductHandler) GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Retrieve the priceGt from the url */
		priceGt, err := strconv.ParseFloat(r.URL.Query().Get("priceGt"), 64)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid priceGt.")
			return
		}

		filteredProducts := p.ProductService.GetProductByPriceGt(priceGt)
		response.JSON(w, http.StatusOK, filteredProducts)
	}
}

// AddNewProduct creates a new product on the website
// URL params : none
// Body params: BodyRequestProductJSON
func (p *ProductHandler) AddNewProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Check if the user is authenthicated */
		if !p.isAuthenthicated(r) {
			response.Text(w, http.StatusUnauthorized, "Unauthorized.")
			return
		}

		/* Retrieve the body from the request */
		var body BodyRequestProductJSON
		err := request.JSON(r, &body)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid body.")
			return
		}

		/* Serialize to internal.TProduct */
		product := internal.TProduct{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		/* Intert the new product into repository */
		err = p.ProductService.InsertNewProduct(&product)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductAlreadyExists):
				response.Text(w, http.StatusBadRequest, "Product already exists.")
				return
			case errors.Is(err, internal.ErrEmptyField), errors.Is(err, internal.ErrInvalidDate):
				response.Text(w, http.StatusBadRequest, "Invalid body."+err.Error())
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}

		/* Serialize to ProductJSON */
		productJSON := ProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		/* Send the new product as response */
		response.JSON(w, http.StatusCreated, map[string]any{
			"data":    productJSON,
			"message": "Product created successfully.",
		})
	}
}

// isProductBodyComplete checks if the body request is complete
// Args:
//
//	fields: map[string]any
//
// Return:
//
//	bool: true if the body is complete, false otherwise
func isProductBodyComplete(fields map[string]any) bool {
	for _, value := range productFields {
		if _, ok := fields[value]; !ok {
			return false
		}
	}
	return true

}

// UpdateProduct update a product on the website
// URL params : none
// Body params: ProductJSON
func (p *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Check if the user is authenthicated */
		if !p.isAuthenthicated(r) {
			response.Text(w, http.StatusUnauthorized, "Unauthorized.")
			return
		}

		/* Check if the request body is complete */
		var fields map[string]any
		err := request.JSON(r, &fields)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid body.")
			return
		}
		if !isProductBodyComplete(fields) {
			response.Text(w, http.StatusBadRequest, "Invalid body.")
			return
		}

		/* Serialize to internal.TProduct */
		product := internal.TProduct{
			ID:          int(fields["id"].(float64)),
			Name:        fields["name"].(string),
			Quantity:    int(fields["quantity"].(float64)),
			CodeValue:   fields["code_value"].(string),
			IsPublished: fields["is_published"].(bool),
			Expiration:  fields["expiration"].(string),
			Price:       fields["price"].(float64),
		}

		/* Update the product into repository */
		err = p.ProductService.UpdateProduct(&product)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductAlreadyExists):
				response.Text(w, http.StatusBadRequest, "Product code already exists.")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}

		/* Send the response to the client */
		response.JSON(w, http.StatusCreated, map[string]any{
			"data":    product,
			"message": "Product updated successfully.",
		})
	}
}

// UpdateProduct partially updates a product on the website
// URL params : id
// Body params: BodyRequestProductJSON
func (p *ProductHandler) UpdateProductPartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Check if the user is authenthicated */
		if !p.isAuthenthicated(r) {
			response.Text(w, http.StatusUnauthorized, "Unauthorized.")
			return
		}

		/* Retrieve the id from the url */
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid ID.")
			return
		}

		/* Retrieve the fields from the request body */
		var fields map[string]any
		err = request.JSON(r, &fields)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid body.")
			return
		}

		/* Retrieve the Product by ID */
		product, err := p.ProductService.GetProductByID(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotExists):
				response.Text(w, http.StatusNotFound, "Product not found.")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}

		/* Update the product fields */
		for key, value := range fields {
			switch key {
			case "name":
				product.Name = value.(string)
			case "quantity":
				product.Quantity = int(value.(float64))
			case "code_value":
				product.CodeValue = value.(string)
			case "is_published":
				product.IsPublished = value.(bool)
			case "expiration":
				product.Expiration = value.(string)
			case "price":
				product.Price = value.(float64)
			default:
				response.Text(w, http.StatusBadRequest, "Invalid body unpespected field: "+key)
				return
			}
		}

		/* Update the product */
		err = p.ProductService.UpdateProduct(&product)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotExists):
				response.Text(w, http.StatusNotFound, "Product not found.")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}

		/* Send the response to the client */
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product updated successfully.",
		})
	}
}

// DeleteProduct deletes a product on the website
// URL params : id
// Body params: ProductJSON
func (p *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Check if the user is authenthicated */
		if !p.isAuthenthicated(r) {
			response.Text(w, http.StatusUnauthorized, "Unauthorized.")
			return
		}

		/* Retrieve the id from the url */
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid ID.")
			return
		}
		/* Delete the product by id */
		err = p.ProductService.DeleteProduct(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotExists):
				response.Text(w, http.StatusNotFound, "Product not found.")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "Internal server error.")
				return
			}
		}
		/* Send the response to the client */
		response.Text(w, http.StatusOK, "Product deleted successfully.")
	}
}
