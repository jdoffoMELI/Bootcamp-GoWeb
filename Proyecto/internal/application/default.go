package application

import (
	"net/http"
	"proyecto/internal/handlers"
	"proyecto/internal/repository"
	"proyecto/internal/service"
	"proyecto/internal/storage"
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

// Run runs the application
func (h *ApplicationDefault) Run() {
	/* Intialize dependencies */
	storage := storage.NewProductStorageDefault("/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Proyecto/docs/db/products.json")
	repository := repository.NewProductMap(storage)
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

		// PUT handlers
		r.Put("/", handler.UpdateProduct())

		// PATCH handlers
		r.Patch("/{id}", handler.UpdateProductPartial())

		// DELETE handlers
		r.Delete("/{id}", handler.DeleteProduct())
	})

	/* Intialize the server */
	http.ListenAndServe(h.address, router)
}
