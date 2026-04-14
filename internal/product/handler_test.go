package product

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeService struct {
	createFn func(ctx context.Context, req CreateProductRequest) (Product, error)
	listFn   func(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error)
}

func (f fakeService) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	return f.createFn(ctx, req)
}

func (f fakeService) List(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error) {
	return f.listFn(ctx, rawQuery, q)
}

func TestHandler_CreateProduct(t *testing.T) {
	h := NewHandler(fakeService{
		createFn: func(ctx context.Context, req CreateProductRequest) (Product, error) {
			return Product{ID: "1", Name: req.Name, Type: req.Type, Price: req.Price}, nil
		},
		listFn: func(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error) {
			return &ListProductsResponse{}, nil
		},
	})

	body, _ := json.Marshal(CreateProductRequest{Name: "Bayam", Type: TypeSayuran, Price: 12000})
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateProduct(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandler_ListProducts_BadQuery(t *testing.T) {
	h := NewHandler(fakeService{
		createFn: func(ctx context.Context, req CreateProductRequest) (Product, error) { return Product{}, nil },
		listFn:   func(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error) { return &ListProductsResponse{}, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/product?type=Invalid", nil)
	rec := httptest.NewRecorder()
	h.ListProducts(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

