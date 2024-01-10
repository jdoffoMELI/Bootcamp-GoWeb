package Product

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
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
//		jsonPath: Json file path.
// Return:
//		[]Product: Slice of products retrieved from a json file.
//		error: 	   Error raised during the execution (if exists).

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

// HasEmptyValues checks if the product has empty values
// HasEmptyValues() -> bool
func (p *TProduct) HasEmptyValues() bool {
	result := p.Name == "" || p.Quantity == 0 || p.CodeValue == "" || p.Expiration == "" || p.Price == 0.0
	return result
}

// HasValidDate checks if the product has a valid date
// HasValidDate() -> bool
func (p *TProduct) HasValidDate() bool {
	tokenSlice := strings.Split(p.Expiration, "/")
	if len(tokenSlice) != 3 {
		return false
	}
	tokenDay, tokenMonth, tokenYear := tokenSlice[0], tokenSlice[1], tokenSlice[2]
	day, err := strconv.Atoi(tokenDay)
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(tokenMonth)
	if err != nil {
		return false
	}
	year, err := strconv.Atoi(tokenYear)
	if err != nil {
		return false
	}
	return day > 0 && day <= 31 && month > 0 && month <= 12 && year > 1900 && year <= 2024
}
