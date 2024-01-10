package server

import (
	"net/http"
	handlers "proyecto/pkg/Handlers"

	"github.com/go-chi/chi/v5"
)

/* ping endpoint routines */
// [GET] ping returns "pong" as response
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ok", "data":"pong"}`))

}

// InitServer initializes all the resources needed by the server and put it on listen mode.
func InitServer() {
	productHandler := handlers.NewProductHandler() // Product handler instantiations
	productHandler.InitHandler()                   // dumps products.json into memory

	router := chi.NewRouter()
	/* Ping endpoints */
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", ping)
	})
	/* Products endpoints */
	router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAllProducts())
		r.Get("/{id}", productHandler.GetProductByID())
		r.Get("/search", productHandler.GetProductByPrice())
	})

	http.ListenAndServe(":8080", router)
}
