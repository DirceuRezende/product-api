package product

import "errors"

var ErrNotFound = errors.New("produto não encontrado")

type ProductRepository interface {
	Save(product Product) (*Product, error)
	Update(product Product) (*Product, error)
	Delete(id uint64) (bool, error)
	FindOne(id uint64) (*Product, error)
	FindAll() ([]*Product, error)
}
