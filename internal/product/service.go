package product

import (
	"context"
	"time"
)

// Input reúne os dados editáveis de um produto, independente do formato de transporte (HTTP).
type Input struct {
	Name        string
	Description string
	Price       uint64
	Stock       uint64
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id uint64) (*Product, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}
	return s.repo.FindOne(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) Create(ctx context.Context, in Input) (*Product, error) {
	now := time.Now()
	return s.repo.Save(ctx, Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (s *Service) Update(ctx context.Context, id uint64, in Input) (*Product, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	p, err := s.repo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Name = in.Name
	p.Description = in.Description
	p.Price = in.Price
	p.Stock = in.Stock
	p.UpdatedAt = time.Now()

	return s.repo.Update(ctx, *p)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
