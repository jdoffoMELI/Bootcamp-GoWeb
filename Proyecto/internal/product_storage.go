package internal

import "errors"

/* Error definition */
var (
	ErrBadFile = errors.New("bad file")
)

/* Product storage definition */
type ProductStorage interface {
	GetAll() (map[int]TProduct, error) // Get all products from storage
	WriteAll(map[int]TProduct) error   // Write all products to storage
}
