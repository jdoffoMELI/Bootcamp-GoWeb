package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"proyecto/internal"
	"proyecto/internal/handlers"
	"proyecto/internal/repository"
	"proyecto/internal/service"
	"proyecto/platform/web/response"

	"github.com/go-chi/chi/v5"
)

type ApplicationDefault struct {
	address string // Server address (host:port)
}

// NewApplicationDefault creates a new default valued ApplicationDefault
// NewApplicationDefault(string) -> *ApplicationDefault
// Args:
//		address: Server address (host:port)
// Return:
//		*ApplicationDefault: New ApplicationDefault instance

func NewApplicationDefault(address string) *ApplicationDefault {
	return &ApplicationDefault{
		address: address,
	}
}

/* ping endpoint routines */
// [GET] ping returns "pong" as response
func ping(w http.ResponseWriter, r *http.Request) {
	response.Text(w, http.StatusOK, "pong")
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

// Run runs the application
func (h *ApplicationDefault) Run() {
	/* Dump the data into memory */
	data, err := dumpJson("/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Proyecto/docs/db/products.json")
	if err != nil {
		panic(err)
	}

	/* Intialize dependencies */
	repository := repository.NewProductMap(sliceToMap(data), len(data))
	service := service.NewProductServiceDefault(repository)
	handler := handlers.NewProductHandler(service)
	router := chi.NewRouter()

	/* Ping endpoints */
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", ping)
	})

	/* Product endpoints */
	router.Route("/products", func(r chi.Router) {
		// GET handlers
		r.Get("/", handler.GetAllProducts())
		r.Get("/{id}", handler.GetProductByID())
		r.Get("/search", handler.GetProductByPrice())

		// POST handlers
		r.Post("/", handler.AddNewProduct())
	})

	/* Intialize the server */
	http.ListenAndServe(h.address, router)
}