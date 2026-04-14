package product

import (
	"context"
	"strings"
)

type Service interface {
	Create(ctx context.Context, req CreateProductRequest) (Product, error)
	List(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error)
}

type service struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) Service {
	return &service{repo: repo, cache: cache}
}

func (s *service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len(req.Name) > 200 {
		return Product{}, ErrInvalidName
	}
	if !req.Type.Valid() {
		return Product{}, ErrInvalidType
	}
	if req.Price <= 0 {
		return Product{}, ErrInvalidPrice
	}
	p, err := s.repo.Create(ctx, req)
	if err != nil {
		return Product{}, err
	}
	_ = s.cache.InvalidateLists(ctx)
	return p, nil
}

func (s *service) List(ctx context.Context, rawQuery string, q ListQuery) (*ListProductsResponse, error) {
	ck := ListCacheKey(rawQuery)
	if v, ok := s.cache.GetList(ctx, ck); ok {
		return v, nil
	}

	items, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, err
	}
	resp := &ListProductsResponse{
		Data:   items,
		Limit:  q.Limit,
		Offset: q.Offset,
		Total:  total,
	}
	_ = s.cache.SetList(ctx, ck, resp)
	return resp, nil
}

