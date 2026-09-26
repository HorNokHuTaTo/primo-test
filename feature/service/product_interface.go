package productservice

import (
	"context"
	productmodel "exam/feature/model"
	productrepository "exam/feature/repository"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req productmodel.CreateProductRequest) (*productmodel.Product, error)
	UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error
}

type ProductService struct {
	ProductRepository productrepository.IProductRepository
}

func NewProductService(productRepository productrepository.IProductRepository) IProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}
