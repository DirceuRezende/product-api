package product

import (
	"errors"
	product "products-api/internal/domain/product"
	"time"
)

type ProductService struct {
	repository product.ProductRepository
}

type ProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Stock       uint64 `json:"stock"`
}

type ProductResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Stock       uint64 `json:"stock"`
}

func NewProductService(repository product.ProductRepository) *ProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetProduct(id uint64) (*ProductResponse, error) {
	if id == 0 {
		return nil, errors.New("ID inválido")
	}

	product, err := p.repository.FindOne(id)

	if err != nil {
		return nil, err
	}

	return &ProductResponse{product.ID, product.Name, product.Description, product.Price, product.Stock}, nil
}

func (p *ProductService) GetProducts() ([]*ProductResponse, error) {
	products, err := p.repository.FindAll()

	var list []*ProductResponse = []*ProductResponse{}

	if err != nil {
		return nil, err
	}

	for _, product := range products {
		list = append(list, &ProductResponse{product.ID, product.Name, product.Description, product.Price, product.Stock})
	}

	return list, nil
}

func (p *ProductService) CreateProduct(request ProductRequest) (*ProductResponse, error) {
	product, err := p.repository.Save(product.Product{ID: 0, Name: request.Name, Description: request.Description, Price: request.Price, Stock: request.Stock, CreateAt: time.Now(), UpdateAt: time.Now()})

	if err != nil {
		return nil, err
	}

	return &ProductResponse{product.ID, product.Name, product.Description, product.Price, product.Stock}, nil
}

func (p *ProductService) AlterProduct(id uint64, request ProductRequest) (*ProductResponse, error) {

	if id == 0 {
		return nil, errors.New("ID inválido")
	}

	product, err := p.repository.FindOne(id)

	if err != nil {
		return nil, err
	}

	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock
	product.UpdateAt = time.Now()

	updatedProduct, err := p.repository.Update(*product)

	if err != nil {
		return nil, err
	}

	return &ProductResponse{updatedProduct.ID, updatedProduct.Name, updatedProduct.Description, updatedProduct.Price, updatedProduct.Stock}, nil
}

func (p *ProductService) DeleteProduct(id uint64) error {
	if id == 0 {
		return errors.New("ID inválido")
	}

	deleted, err := p.repository.Delete(id)

	if err != nil {
		return err
	}

	if !deleted {
		return product.ErrNotFound
	}

	return nil
}
