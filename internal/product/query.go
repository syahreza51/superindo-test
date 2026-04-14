package product

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type SortBy string

const (
	SortByDate  SortBy = "date"
	SortByPrice SortBy = "price"
	SortByName  SortBy = "name"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type ListQuery struct {
	ID    string
	Name  string
	Types []Type

	SortBy    SortBy
	SortOrder SortOrder

	Limit  int
	Offset int
}

func ParseListQuery(r *http.Request) (ListQuery, error) {
	q := r.URL.Query()

	limit := parseIntDefault(q.Get("limit"), 20)
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := parseIntDefault(q.Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}

	sortBy := SortBy(strings.ToLower(strings.TrimSpace(q.Get("sort_by"))))
	if sortBy == "" {
		sortBy = SortByDate
	}
	if sortBy != SortByDate && sortBy != SortByPrice && sortBy != SortByName {
		return ListQuery{}, ErrInvalidSortBy
	}

	sortOrder := SortOrder(strings.ToLower(strings.TrimSpace(q.Get("sort_order"))))
	if sortOrder == "" {
		sortOrder = SortDesc
	}
	if sortOrder != SortAsc && sortOrder != SortDesc {
		sortOrder = SortDesc
	}

	var types []Type
	for _, t := range q["type"] {
		tt := Type(strings.TrimSpace(t))
		if tt == "" {
			continue
		}
		if !tt.Valid() {
			return ListQuery{}, ErrInvalidType
		}
		types = append(types, tt)
	}

	id := strings.TrimSpace(q.Get("id"))
	if id != "" {
		if _, err := uuid.Parse(id); err != nil {
			return ListQuery{}, errors.New("invalid product id")
		}
	}

	return ListQuery{
		ID:        id,
		Name:      strings.TrimSpace(q.Get("name")),
		Types:     types,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}, nil
}

func parseIntDefault(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

