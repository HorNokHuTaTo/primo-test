package productmodel

import (
	"encoding/json"
	"time"
)

// CreateProduct
type (
	CreateProductRequest struct {
		Name        string   `json:"name" validate:"required,max=255"`
		Description *string  `json:"description" validate:"omitempty"`
		SalePrice   *float64 `json:"sale_price" validate:"omitempty"`
		Price       float64  `json:"price" validate:"required"`
	}

	Product struct {
		ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
		Name        string    `gorm:"not null;type:varchar(255)" json:"name"`
		Description *string   `gorm:"type:text" json:"description"`
		SalePrice   *float64  `gorm:"type:numeric(10,2)" json:"sale_price"`
		Price       float64   `gorm:"not null;type:numeric(10,2)" json:"price"`
		CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	}

	ProductResponse struct {
		Successful bool     `json:"successful"`
		ErrorCode  string   `json:"error_code,omitempty"`
		Data       *Product `json:"data,omitempty"`
	}
)

// UpdateProduct
type (
	UpdateProductRequest struct {
		Name        *string              `json:"name" validate:"omitempty,max=255"`
		Description FieldUpdate[string]  `json:"description" validate:"omitempty" swaggertype:"string"`
		SalePrice   FieldUpdate[float64] `json:"sale_price" validate:"omitempty" swaggertype:"float64"`
		Price       *float64             `json:"price" validate:"omitempty"`
	}
	FieldUpdate[T any] struct {
		Present bool
		Value   *T
	}
)

func (f *FieldUpdate[T]) UnmarshalJSON(data []byte) error {
	f.Present = true
	if string(data) == "null" {
		f.Value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f.Value = &v
	return nil
}
