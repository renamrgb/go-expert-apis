package entity

import (
	"errors"
	"time"

	"github.com/renamrgb/go-expert-apis/pkg/entity"
)

var (
	ErrIDIsRequired  = errors.New("id is required")
	ErrInvalidID     = errors.New("invalid id")
	ErrNameRequired  = errors.New("name is required")
	ErrPriceRequired = errors.New("price  is required")
	ErrInvalidPrice  = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price int) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price == 0 {
		return ErrInvalidPrice
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
