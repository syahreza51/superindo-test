package product

type CreateProductRequest struct {
	Name  string `json:"name"`
	Type  Type   `json:"type"`
	Price int64  `json:"price"`
}

type ListProductsResponse struct {
	Data   []Product `json:"data"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	Total  int       `json:"total"`
}

