package product

import (
	"context"
	"testing"
)

type fakeRepo struct {
	created Product
}

func (f *fakeRepo) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	f.created = Product{ID: "id", Name: req.Name, Type: req.Type, Price: req.Price}
	return f.created, nil
}

func (f *fakeRepo) List(ctx context.Context, q ListQuery) ([]Product, int, error) {
	return []Product{{ID: "1", Name: "Bayam", Type: TypeSayuran, Price: 12000}}, 1, nil
}

type noopCache struct{}

func (n noopCache) GetList(ctx context.Context, cacheKey string) (*ListProductsResponse, bool) { return nil, false }
func (n noopCache) SetList(ctx context.Context, cacheKey string, v *ListProductsResponse) error {
	return nil
}
func (n noopCache) InvalidateLists(ctx context.Context) error { return nil }

func TestService_CreateValidation(t *testing.T) {
	svc := NewService(&fakeRepo{}, noopCache{})

	if _, err := svc.Create(context.Background(), CreateProductRequest{Name: "", Type: TypeSayuran, Price: 1}); err == nil {
		t.Fatalf("expected error for empty name")
	}
	if _, err := svc.Create(context.Background(), CreateProductRequest{Name: "A", Type: Type("Invalid"), Price: 1}); err == nil {
		t.Fatalf("expected error for invalid type")
	}
	if _, err := svc.Create(context.Background(), CreateProductRequest{Name: "A", Type: TypeSayuran, Price: 0}); err == nil {
		t.Fatalf("expected error for invalid price")
	}
	if _, err := svc.Create(context.Background(), CreateProductRequest{Name: "Bayam", Type: TypeSayuran, Price: 12000}); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

