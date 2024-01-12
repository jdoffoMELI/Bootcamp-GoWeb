package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// MiddlewareLogger is a middleware that logs the requests to the server
// MiddlewareLogger(*os.File) -> func(http.Handler) http.Handler
// Args:
//		file: File where the log will be written
// Return:
//		func(http.Handler) http.Handler: HTTP handler function

func MiddlewareLogger(file *os.File) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* Writes the log message */
			logString := fmt.Sprintf(
				"[%s] %s %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				r.Method,
				r.URL.Path)
			_, err := file.WriteString(logString)
			if err != nil {
				panic(err)
			}

			/* Calls the next handler */
			handler.ServeHTTP(w, r)
		})
	}
}
