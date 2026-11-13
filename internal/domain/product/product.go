package product

import "time"

type Product struct {
	ID          uint64
	Name        string
	Description string
	Price       uint64
	Stock       uint64
	CreateAt    time.Time
	UpdateAt    time.Time
}
