package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/renamrgb/go-expert-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const idQuery = "id = ?"

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "falha ao abrir banco in-memory")
	require.NoError(t, db.AutoMigrate(&entity.Product{}))
	return db
}

func newProductDB(t *testing.T) (*Product, *gorm.DB) {
	t.Helper()
	db := setupDB(t)
	return NewProduct(db), db
}

func seedProducts(t *testing.T, db *gorm.DB, n int) []*entity.Product {
	t.Helper()
	products := make([]*entity.Product, n)
	for i := range products {
		p, err := entity.NewProduct(fmt.Sprintf("Product %d", i+1), rand.Float64()*100)
		require.NoError(t, err)
		require.NoError(t, db.Create(p).Error)
		products[i] = p
	}
	return products
}

func TestCreateProduct(t *testing.T) {
	productDB, db := newProductDB(t)
	product, err := entity.NewProduct("Product 1", 10.0)
	require.NoError(t, err)

	require.NoError(t, productDB.Create(product))

	var found entity.Product
	require.NoError(t, db.First(&found, idQuery, product.ID).Error)
	assert.Equal(t, product.ID, found.ID)
	assert.Equal(t, product.Name, found.Name)
	assert.Equal(t, product.Price, found.Price)
}

func TestFindByID(t *testing.T) {
	productDB, db := newProductDB(t)
	products := seedProducts(t, db, 1)
	product := products[0]

	t.Run("produto existente", func(t *testing.T) {
		found, err := productDB.FindByID(product.ID.String())
		require.NoError(t, err)
		assert.Equal(t, product.ID, found.ID)
		assert.Equal(t, product.Name, found.Name)
		assert.Equal(t, product.Price, found.Price)
	})

	t.Run("id inexistente retorna erro", func(t *testing.T) {
		_, err := productDB.FindByID("00000000-0000-0000-0000-000000000000")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestUpdateProduct(t *testing.T) {
	productDB, db := newProductDB(t)
	products := seedProducts(t, db, 1)
	product := products[0]

	product.Name = "Updated Product"
	product.Price = 20.0
	require.NoError(t, productDB.Update(product))

	var found entity.Product
	require.NoError(t, db.First(&found, idQuery, product.ID).Error)
	assert.Equal(t, product.ID, found.ID)
	assert.Equal(t, "Updated Product", found.Name)
	assert.Equal(t, 20.0, found.Price)
}

func TestDeleteProduct(t *testing.T) {
	productDB, db := newProductDB(t)
	products := seedProducts(t, db, 1)
	product := products[0]

	require.NoError(t, productDB.Delete(product.ID.String()))

	var found entity.Product
	err := db.First(&found, idQuery, product.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFindAllProducts(t *testing.T) {
	productDB, db := newProductDB(t)
	seedProducts(t, db, 23)

	tests := []struct {
		name          string
		page, limit   int
		sort          string
		wantLen       int
		wantFirstName string
		wantLastName  string
	}{
		{
			name: "página 1 — 10 itens",
			page: 1, limit: 10, sort: "asc",
			wantLen:       10,
			wantFirstName: "Product 1",
			wantLastName:  "Product 10",
		},
		{
			name: "página 2 — 10 itens",
			page: 2, limit: 10, sort: "asc",
			wantLen:       10,
			wantFirstName: "Product 11",
			wantLastName:  "Product 20",
		},
		{
			name: "página 3 — 3 itens restantes",
			page: 3, limit: 10, sort: "asc",
			wantLen:       3,
			wantFirstName: "Product 21",
			wantLastName:  "Product 23",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			products, err := productDB.FindAll(tc.page, tc.limit, tc.sort)
			require.NoError(t, err)
			assert.Len(t, products, tc.wantLen)
			assert.Equal(t, tc.wantFirstName, products[0].Name)
			assert.Equal(t, tc.wantLastName, products[tc.wantLen-1].Name)
		})
	}
}
