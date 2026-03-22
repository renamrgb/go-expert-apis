package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgentity "github.com/renamrgb/go-expert-apis/pkg/entity"
)

func TestNewProductSuccess(t *testing.T) {
	product, err := NewProduct("Product 1", 100.00)

	require.NoError(t, err)
	require.NotNil(t, product)

	assert.NotEmpty(t, product.ID.String())
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, 100.00, product.Price)
	assert.WithinDuration(t, time.Now(), product.CreatedAt, time.Second)
}

func TestNewProductErrorScenarios(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		inputPrice  float64
		expectedErr error
	}{
		{
			name:        "name empty",
			inputName:   "",
			inputPrice:  100,
			expectedErr: ErrNameRequired,
		},
		{
			name:        "price zero",
			inputName:   "Product",
			inputPrice:  0,
			expectedErr: ErrInvalidPrice,
		},
		{
			name:        "price negative",
			inputName:   "Product",
			inputPrice:  -10,
			expectedErr: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.inputName, tt.inputPrice)

			require.Error(t, err)
			assert.Nil(t, product)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestProductValidateSuccess(t *testing.T) {
	id := pkgentity.NewID()

	product := &Product{
		ID:        id,
		Name:      "Valid Product",
		Price:     100,
		CreatedAt: time.Now(),
	}

	err := product.Validate()

	require.NoError(t, err)
}

func TestProductValidateErrorScenarios(t *testing.T) {
	validID := pkgentity.NewID()

	tests := []struct {
		name        string
		product     Product
		expectedErr error
	}{
		{
			name: "name empty",
			product: Product{
				ID:    validID,
				Name:  "",
				Price: 100,
			},
			expectedErr: ErrNameRequired,
		},
		{
			name: "price zero",
			product: Product{
				ID:    validID,
				Name:  "Product",
				Price: 0,
			},
			expectedErr: ErrInvalidPrice,
		},
		{
			name: "price negative",
			product: Product{
				ID:    validID,
				Name:  "Product",
				Price: -1,
			},
			expectedErr: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()

			require.Error(t, err)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
