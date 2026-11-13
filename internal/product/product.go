package product

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound  = errors.New("produto não encontrado")
	ErrInvalidID = errors.New("ID inválido")
)

type Product struct {
	ID          uint64
	Name        string
	Description string
	Price       uint64
	Stock       uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository é consumido pelo Service; as implementações ficam em subpacotes (ex.: postgres).
type Repository interface {
	Save(ctx context.Context, p Product) (*Product, error)
	Update(ctx context.Context, p Product) (*Product, error)
	Delete(ctx context.Context, id uint64) error
	FindOne(ctx context.Context, id uint64) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
}
