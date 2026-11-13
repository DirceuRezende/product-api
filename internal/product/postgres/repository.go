package postgres

import (
	"context"
	"errors"
	"fmt"

	"products-api/internal/product"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, p product.Product) (*product.Product, error) {
	const query = `
		INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, price, stock, created_at, updated_at`

	var saved product.Product
	err := r.db.QueryRow(ctx, query,
		p.Name, p.Description, p.Price, p.Stock, p.CreatedAt, p.UpdatedAt,
	).Scan(&saved.ID, &saved.Name, &saved.Description, &saved.Price, &saved.Stock, &saved.CreatedAt, &saved.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return &saved, nil
}

func (r *Repository) Update(ctx context.Context, p product.Product) (*product.Product, error) {
	const query = `
		UPDATE products
		SET name = $2, description = $3, price = $4, stock = $5, updated_at = $6
		WHERE id = $1
		RETURNING id, name, description, price, stock, created_at, updated_at`

	var updated product.Product
	err := r.db.QueryRow(ctx, query,
		p.ID, p.Name, p.Description, p.Price, p.Stock, p.UpdatedAt,
	).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Price, &updated.Stock, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return product.ErrNotFound
	}

	return nil
}

func (r *Repository) FindOne(ctx context.Context, id uint64) (*product.Product, error) {
	const query = `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1`

	var found product.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&found.ID, &found.Name, &found.Description, &found.Price, &found.Stock, &found.CreatedAt, &found.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("find product: %w", err)
	}

	return &found, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*product.Product, error) {
	const query = `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	defer rows.Close()

	products := []*product.Product{}

	for rows.Next() {
		var p product.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}

	return products, nil
}
