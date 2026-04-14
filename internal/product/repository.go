package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, req CreateProductRequest) (Product, error)
	List(ctx context.Context, q ListQuery) ([]Product, int, error)
}

type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO products (name, type, price)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, type, price, created_at`,
		req.Name, string(req.Type), req.Price,
	)

	var p Product
	var typ string
	if err := row.Scan(&p.ID, &p.Name, &typ, &p.Price, &p.CreatedAt); err != nil {
		return Product{}, err
	}
	p.Type = Type(typ)
	return p, nil
}

func (r *pgRepository) List(ctx context.Context, q ListQuery) ([]Product, int, error) {
	where, args := buildWhere(q)

	countSQL := "SELECT count(1) FROM products" + where
	var total int
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := buildOrderBy(q)
	listSQL := `
SELECT id, name, type, price, created_at
FROM products` + where + orderBy + ` LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, q.Limit, q.Offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Product, error) {
		var p Product
		var typ string
		if err := row.Scan(&p.ID, &p.Name, &typ, &p.Price, &p.CreatedAt); err != nil {
			return Product{}, err
		}
		p.Type = Type(typ)
		return p, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func buildWhere(q ListQuery) (string, []any) {
	var (
		clauses []string
		args    []any
	)
	if q.ID != "" {
		args = append(args, q.ID)
		clauses = append(clauses, fmt.Sprintf("id = $%d", len(args)))
	}
	if q.Name != "" {
		args = append(args, "%"+q.Name+"%")
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if len(q.Types) > 0 {
		placeholders := make([]string, 0, len(q.Types))
		for _, t := range q.Types {
			args = append(args, string(t))
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}
		clauses = append(clauses, "type IN ("+strings.Join(placeholders, ",")+")")
	}
	if len(clauses) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func buildOrderBy(q ListQuery) string {
	col := "created_at"
	switch q.SortBy {
	case SortByPrice:
		col = "price"
	case SortByName:
		col = "name"
	case SortByDate:
		col = "created_at"
	}
	dir := "DESC"
	if q.SortOrder == SortAsc {
		dir = "ASC"
	}
	if col == "name" {
		return " ORDER BY name " + dir + ", created_at DESC"
	}
	return " ORDER BY " + col + " " + dir
}

