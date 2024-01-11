package response

import "net/http"

// Text writes a text response to the client.
// Text(w http.ResponseWriter, code int, body string)
// Args:
//		w    :  HTTP response writer.
//		code :  HTTP status code.
//		body :  Response body.
// Return:
//		none

func Text(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(body))
}
