package application

import (
	"net/http"
	"os"
	"proyecto/internal/handlers"
	"proyecto/internal/middleware"
	"proyecto/internal/repository"
	"proyecto/internal/service"
	"proyecto/internal/storage"

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

// Run runs the application
func (h *ApplicationDefault) Run() {
	/* Intialize dependencies */
	storagePath := "/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Proyecto/docs/db/products.json"
	logpath := "/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Proyecto/docs/logs/log.txt"
	storage := storage.NewProductStorageDefault(storagePath)
	repository := repository.NewProductMap(storage)
	service := service.NewProductServiceDefault(repository)
	handler := handlers.NewProductHandler(service)
	router := chi.NewRouter()
	/* Open log file */
	file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	/* Middlewares */
	router.Use(middleware.MiddlewareLogger(file))
	router.Use(middleware.MiddelwareAuthentication)

	/* Public endpoints */
	router.Route("/products", func(r chi.Router) {
		/*
			TODO:
				No se pueden hacer Route con la misma ruta base por lo que
			de este modo no se puede aplicar middlewares difenciados por tipo de endpoint.
				Intente usar group pero esto siempre retorna 405 Method Not Allowed.

				Una solucion posible es aplicar el middleware (de autenticacion) directamente
			sobre las funciones handler de los endpoins privados, pero esto no es muy elegante.
		*/

		/* Public Endpoints */
		r.Get("/", handler.GetAllProducts())
		r.Get("/{id}", handler.GetProductByID())
		r.Get("/search", handler.GetProductByPrice())

		/* Private Endpoints */
		r.Post("/", handler.AddNewProduct())
		r.Put("/", handler.UpdateProduct())
		r.Patch("/{id}", handler.UpdateProductPartial())
		r.Delete("/{id}", handler.DeleteProduct())
	})

	http.ListenAndServe(h.address, router)
}
