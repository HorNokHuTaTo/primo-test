package productservice

import (
	"context"
	productmodel "exam/feature/model"
)

func (s *ProductService) CreateProduct(ctx context.Context, req productmodel.CreateProductRequest) (*productmodel.Product, error) {
	createProduct := productmodel.Product{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	res, err := s.ProductRepository.CreateProduct(ctx, createProduct)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error {
	if err := s.ProductRepository.UpdateProduct(ctx, productID, req); err != nil {
		return err
	}

	return nil
}
