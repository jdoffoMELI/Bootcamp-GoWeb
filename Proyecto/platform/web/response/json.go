package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response to the client.
// JSON(w http.ResponseWriter, code int, body any)
// Args:
//		w    :  HTTP response writer.
//		code :  HTTP status code.
//		body :  Response body.
// Return:
//		none

func JSON(w http.ResponseWriter, code int, body any) {
	/* Checks empty body */
	if body == nil {
		w.WriteHeader(code)
		return
	}

	/* Byte encoding the body */
	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/* Writes the response */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}
