package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/* Errors definition */
var (
	ErrRequestContentTypeNotJSON = errors.New("request content type is not application/json")
	ErrRequestJSONInvalid        = errors.New("request json invalid")
)

// JSON decodes the request body into the given pointer.
// JSON(r *http.Request, ptr any) -> (err :error)
// Args:
//		r    :	HTTP request to decode.
//		ptr  :  Target data structure to decode into.
// Return:
//		err  :  Error raised during the execution (if exists).

func JSON(r *http.Request, ptr any) (err error) {
	// Checks if the request content type is application/json
	if r.Header.Get("Content-Type") != "application/json" {
		err = ErrRequestContentTypeNotJSON
		return
	}
	// Decodes the request body into data structure
	err = json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrRequestJSONInvalid, err)
		return
	}
	return
}
