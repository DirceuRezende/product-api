package postgres

import (
	"context"
	"errors"
	"fmt"

	product "products-api/internal/domain/product"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (p *PostgresProductRepository) Save(prod product.Product) (*product.Product, error) {
	const query = `
		INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, price, stock, created_at, updated_at`

	var saved product.Product
	err := p.db.QueryRow(context.Background(), query,
		prod.Name, prod.Description, prod.Price, prod.Stock, prod.CreateAt, prod.UpdateAt,
	).Scan(&saved.ID, &saved.Name, &saved.Description, &saved.Price, &saved.Stock, &saved.CreateAt, &saved.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return &saved, nil
}

func (p *PostgresProductRepository) Update(prod product.Product) (*product.Product, error) {
	const query = `
		UPDATE products
		SET name = $2, description = $3, price = $4, stock = $5, updated_at = $6
		WHERE id = $1
		RETURNING id, name, description, price, stock, created_at, updated_at`

	var updated product.Product
	err := p.db.QueryRow(context.Background(), query,
		prod.ID, prod.Name, prod.Description, prod.Price, prod.Stock, prod.UpdateAt,
	).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Price, &updated.Stock, &updated.CreateAt, &updated.UpdateAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &updated, nil
}

func (p *PostgresProductRepository) Delete(id uint64) (bool, error) {
	tag, err := p.db.Exec(context.Background(), `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("delete product: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (p *PostgresProductRepository) FindOne(id uint64) (*product.Product, error) {
	const query = `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1`

	var found product.Product
	err := p.db.QueryRow(context.Background(), query, id).Scan(
		&found.ID, &found.Name, &found.Description, &found.Price, &found.Stock, &found.CreateAt, &found.UpdateAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("find product: %w", err)
	}

	return &found, nil
}

func (p *PostgresProductRepository) FindAll() ([]*product.Product, error) {
	const query = `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY id`

	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	defer rows.Close()

	products := []*product.Product{}

	for rows.Next() {
		var prod product.Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CreateAt, &prod.UpdateAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, &prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}

	return products, nil
}
