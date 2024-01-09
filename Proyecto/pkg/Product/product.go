package Product

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

// TProduct representens a product on the website.
type TProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

// DumpJson creates a slice of products from a json file
// DumpJson(string) -> ([]TProduct, error)
// Args:
//
//	jsonPath: Json file path.
//
// Return:
//
//	[]Product: Slice of products retrieved from a json file.
//	error: 	   Error raised during the execution (if exists).
func DumpJson(jsonPath string) ([]TProduct, error) {
	var jsonSlice []TProduct
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
