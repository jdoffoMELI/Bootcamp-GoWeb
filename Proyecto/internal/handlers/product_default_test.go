package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"proyecto/internal"
	"proyecto/internal/handlers"
	"proyecto/internal/middleware"
	"proyecto/internal/repository"
	"proyecto/internal/service"
	"proyecto/internal/storage"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// initStorage initializes the storage
// initStorage(map[int]internal.TProduct) -> storage.ProductStorageDefault
// Args:
// 	initialProducts: Initial products
// Returns:
// 	ProductStorageDefault: Initialized storage

func initStorage(initialProducts map[int]internal.TProduct) storage.ProductStorageDefault {
	/* Storage creation */
	filepath := "/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Proyecto/internal/handlers/products_test.json"
	storage := storage.NewProductStorageDefault(filepath)

	/* Initial data of the storage */
	err := storage.WriteAll(initialProducts)
	if err != nil {
		panic(err)
	}
	return *storage
}

// addURLParams adds the URL params to the request (needed by Chi framework)
// addURLParams(*http.Request, map[string]string) -> *http.Request
// Args:
// 	req: Request
// 	params: URL params
// Returns:
// 	*http.Request: Request with the URL params

func addURLParams(req *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for key, value := range params {
		chiCtx.URLParams.Add(key, value)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
}

// TestGetAllProducts test the GetAllProducts handler
func TestGetAllProducts(t *testing.T) {
	// Test 1: should return a list of products
	t.Run("should return a list of products", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
			3: {ID: 3, Name: "Product 3", Quantity: 30, CodeValue: "AX03", IsPublished: false, Expiration: "11/11/2003", Price: 30.5},
			4: {ID: 4, Name: "Product 4", Quantity: 40, CodeValue: "AX04", IsPublished: true, Expiration: "11/11/2004", Price: 40.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		handler.GetAllProducts()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusOK
		expectedBody := `{"data": [
			{"id": 1, "name": "Product 1", "quantity": 10, "code_value": "AX01", "is_published": false, "expiration": "11/11/2001", "price": 10.5},
			{"id": 2, "name": "Product 2", "quantity": 20, "code_value": "AX02", "is_published": true, "expiration": "11/11/2002", "price": 20.5},
			{"id": 3, "name": "Product 3", "quantity": 30, "code_value": "AX03", "is_published": false, "expiration": "11/11/2003", "price": 30.5},
			{"id": 4, "name": "Product 4", "quantity": 40, "code_value": "AX04", "is_published": true, "expiration": "11/11/2004", "price": 40.5}
		]}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}

// TestGetProductById test the GetProductById handler
func TestGetProductById(t *testing.T) {
	// Test 1: should return a product
	t.Run("should return a product", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/2", nil)
		req = addURLParams(req, map[string]string{"id": "2"})
		res := httptest.NewRecorder()

		handler.GetProductByID()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusOK
		expectedBody := `{"data":
			{"id":2, "name": "Product 2", "quantity": 20, "code_value": "AX02", "is_published": true, "expiration": "11/11/2002", "price": 20.5}
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	// Test 2: should return a bad request error
	t.Run("should return a bad request error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/A2", nil)
		req = addURLParams(req, map[string]string{"id": "A2"})
		res := httptest.NewRecorder()

		handler.GetProductByID()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusBadRequest
		expectedBody := "Invalid ID."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	// Test 3: should return a not found error
	t.Run("should return a not found error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/3", nil)
		req = addURLParams(req, map[string]string{"id": "3"})
		res := httptest.NewRecorder()

		handler.GetProductByID()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusNotFound
		expectedBody := "Product not found."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}

// TestAddNewProduct test the AddNewProduct handler
func TestAddNewProduct(t *testing.T) {
	// Test 1: should add a new product
	t.Run("should add a new product", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		reqBody := `{
			"name": "new product",
			"quantity": 1000,
			"is_published": true,
			"code_value": "AX04",
			"expiration": "01/01/2000",
			"price": 20
		}`
		req := httptest.NewRequest("POST", "/products/", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.AddNewProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusCreated
		expectedBody := `
		{
			"data": {
				"id": 3,
				"name": "new product",
				"quantity": 1000,
				"code_value": "AX04",
				"is_published": true,
				"expiration": "01/01/2000",
				"price": 20
			},
			"message": "Product created successfully."
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	// Test 2: should return an not authorized error
	t.Run("should return an not authorized error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Set environment variables */
		os.Setenv("TOKEN", "wrong code") // Token to access data modification operations
		defer func() { os.Setenv("TOKEN", "") }()

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = addURLParams(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()

		middleware.MiddelwareAuthentication(handler.AddNewProduct()).ServeHTTP(res, req)

		/* Expected values definition */
		expectedCode := http.StatusUnauthorized
		expectedBody := "Unauthorized."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}

// TestDeleteProduct test the DeleteProduct handler
func TestDeleteProduct(t *testing.T) {
	// Test 1: should delete a product
	t.Run("should delete a product", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("DELETE", "/products/2", nil)
		req = addURLParams(req, map[string]string{"id": "2"})
		res := httptest.NewRecorder()

		handler.DeleteProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusNoContent
		expectedBody := ""

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
	})

	// Test 2: should return a bad request error
	t.Run("should return a bad request error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/A2", nil)
		req = addURLParams(req, map[string]string{"id": "A2"})
		res := httptest.NewRecorder()

		handler.DeleteProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusBadRequest
		expectedBody := "Invalid ID."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	// Test 3: should return a not found error
	t.Run("should return a not found error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/3", nil)
		req = addURLParams(req, map[string]string{"id": "3"})
		res := httptest.NewRecorder()

		handler.DeleteProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusNotFound
		expectedBody := "Product not found."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	// Test 4: should return an not authorized error
	t.Run("should return an not authorized error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Set environment variables */
		os.Setenv("TOKEN", "wrong code") // Token to access data modification operations
		defer func() { os.Setenv("TOKEN", "") }()
		/* Prepare the request and the response */
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = addURLParams(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()

		middleware.MiddelwareAuthentication(handler.DeleteProduct()).ServeHTTP(res, req)

		/* Expected values definition */
		expectedCode := http.StatusUnauthorized
		expectedBody := "Unauthorized."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}

// TestUpdateProduct tests the UpdateProduct handler
func TestUpdateProduct(t *testing.T) {
	// Test 1: should update a product
	t.Run("should update a product", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		reqbody := `{
			"id": 1,
			"name": "new product",
			"quantity": 1000,
			"is_published": true,
			"code_value": "AX04",
			"expiration": "01/01/2000",
			"price": 20
		}`
		req := httptest.NewRequest("PUT", "/products/", strings.NewReader(reqbody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.UpdateProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusOK
		expectedBody := `
		{
			"data": {
				"id": 1,
				"name": "new product",
				"quantity": 1000,
				"code_value": "AX04",
				"is_published": true,
				"expiration": "01/01/2000",
				"price": 20
			},
			"message": "Product updated successfully."
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	// Test 2: should return a bad request error
	t.Run("should return a bad request error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		reqbody := `{
			"id": 1,
			"is_published": true,
			"code_value": "AX04",
			"expiration": "01/01/2000",
			"price": 20
		}`
		req := httptest.NewRequest("PUT", "/products/", strings.NewReader(reqbody))
		res := httptest.NewRecorder()
		handler.UpdateProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusBadRequest
		expectedBody := "Invalid body."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	// Test 3: should return a not found error
	t.Run("should return a not found error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Prepare the request and the response */
		reqbody := `{
			"id": 3,
			"name": "new product",
			"quantity": 1000,
			"is_published": true,
			"code_value": "AX04",
			"expiration": "01/01/2000",
			"price": 20
		}`
		req := httptest.NewRequest("GET", "/products/1", strings.NewReader(reqbody))
		req.Header.Set("Content-Type", "application/json")
		req = addURLParams(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()

		handler.UpdateProduct()(res, req)

		/* Expected values definition */
		expectedCode := http.StatusNotFound
		expectedBody := "Product not found."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	// Test 4: should return an not authorized error
	t.Run("should return an not authorized error", func(t *testing.T) {
		/* Prepare the test data */
		initialProducts := map[int]internal.TProduct{
			1: {ID: 1, Name: "Product 1", Quantity: 10, CodeValue: "AX01", IsPublished: false, Expiration: "11/11/2001", Price: 10.5},
			2: {ID: 2, Name: "Product 2", Quantity: 20, CodeValue: "AX02", IsPublished: true, Expiration: "11/11/2002", Price: 20.5},
		}

		/* Initialize dependencies */
		storage := initStorage(initialProducts)
		repository := repository.NewProductMap(&storage)
		service := service.NewProductServiceDefault(repository)
		handler := handlers.NewProductHandler(service)

		/* Set environment variables */
		os.Setenv("TOKEN", "wrong code") // Token to access data modification operations
		defer func() { os.Setenv("TOKEN", "") }()

		/* Prepare the request and the response */
		reqbody := `{
			"id": 1,
			"name": "new product",
			"quantity": 1000,
			"is_published": true,
			"code_value": "AX04",
			"expiration": "01/01/2000",
			"price": 20
		}`
		req := httptest.NewRequest("GET", "/products/1", strings.NewReader(reqbody))
		req.Header.Set("Content-Type", "application/json")
		req = addURLParams(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()

		middleware.MiddelwareAuthentication(handler.UpdateProduct()).ServeHTTP(res, req)

		/* Expected values definition */
		expectedCode := http.StatusUnauthorized
		expectedBody := "Unauthorized."
		expectedHeader := http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}

		/* Assertions */
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}
