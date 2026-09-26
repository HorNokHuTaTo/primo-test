package productrepository

import (
	"context"
	productmodel "exam/feature/model"
	"time"
)

func (r *ProductRepository) CreateProduct(ctx context.Context, req productmodel.Product) (*productmodel.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.DB.WithContext(ctx).Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	structUpdates := productmodel.Product{}
	if req.Name != nil {
		structUpdates.Name = *req.Name
	}
	if req.Price != nil {
		structUpdates.Price = *req.Price
	}

	if structUpdates.Name != "" || structUpdates.Price != 0 {
		if err := r.DB.WithContext(ctx).
			Model(&productmodel.Product{}).
			Where("id = ?", productID).
			Updates(structUpdates).Error; err != nil {
			return err
		}
	}

	nullableUpdates := map[string]interface{}{}
	if req.Description.Present {
		nullableUpdates["description"] = req.Description.Value
	}
	if req.SalePrice.Present {
		nullableUpdates["sale_price"] = req.SalePrice.Value
	}

	if len(nullableUpdates) > 0 {
		if err := r.DB.WithContext(ctx).
			Model(&productmodel.Product{}).
			Where("id = ?", productID).
			Updates(nullableUpdates).Error; err != nil {
			return err
		}
	}

	return nil
}
