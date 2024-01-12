package middleware

import (
	"net/http"
	"os"
	"proyecto/platform/web/response"
)

// isAuthenthicated returns true if the user is authenthicated using a token
// isAuthenthicated(*http.Request) -> bool
// Args:
//		r: HTTP request
// Return:
//		bool: True if the user is authenthicated, false otherwise

func isAuthenthicated(r *http.Request) bool {
	userToken := r.Header.Get("TOKEN")
	return userToken == os.Getenv("TOKEN")
}

// MiddelwareAuthentication is a middleware that checks if the user is authenthicated using a token
// MiddelwareAuthentication(http.HandlerFunc) -> http.HandlerFunc
// Args:
//		handlerFunc: HTTP handler function
// Return:
//		http.HandlerFunc: HTTP handler function

func MiddelwareAuthentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAuthenthicated(r) {
			handler.ServeHTTP(w, r)
		} else {
			response.Text(w, http.StatusUnauthorized, "Unauthorized.")
		}
	})
}
