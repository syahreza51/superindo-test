package product

import "errors"

var (
	ErrInvalidType   = errors.New("invalid product type")
	ErrInvalidName   = errors.New("invalid product name")
	ErrInvalidPrice  = errors.New("invalid product price")
	ErrInvalidSortBy = errors.New("invalid sort_by")
)

