package productrepository

import (
	"context"
	productmodel "exam/feature/model"

	"gorm.io/gorm"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, req productmodel.Product) (*productmodel.Product, error)
	UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error
}

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{
		DB: db,
	}
}
