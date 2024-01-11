package handlers

import (
	"errors"
	"net/http"
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
//		id (Numeric): ID of the desirable product.

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
//		priceGt (Numeric): Price to filter by.

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
