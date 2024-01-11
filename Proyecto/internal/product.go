package internal

// TProduct representens a product on the website.
type TProduct struct {
	ID          int
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  string
	Price       float64
}
