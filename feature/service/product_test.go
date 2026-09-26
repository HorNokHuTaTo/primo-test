package productservice_test

import (
	"context"
	"errors"
	productmodel "exam/feature/model"
	productservice "exam/feature/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, req productmodel.Product) (*productmodel.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productmodel.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error {
	args := m.Called(ctx, productID, req)
	return args.Error(0)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	desc := "A description"
	salePrice := 79.99
	req := productmodel.CreateProductRequest{
		Name:        "Test Product",
		Description: &desc,
		SalePrice:   &salePrice,
		Price:       99.99,
	}

	expectedProduct := productmodel.Product{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	returnedProduct := &productmodel.Product{
		ID:          "11111111-1111-1111-1111-111111111111",
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	mockRepo.On("CreateProduct", mock.Anything, expectedProduct).Return(returnedProduct, nil)

	result, err := svc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", result.ID)
	assert.Equal(t, "Test Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_MinimalFields_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	req := productmodel.CreateProductRequest{
		Name:  "Minimal Product",
		Price: 10.0,
	}

	expectedProduct := productmodel.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	returnedProduct := &productmodel.Product{
		ID:    "22222222-2222-2222-2222-222222222222",
		Name:  req.Name,
		Price: req.Price,
	}

	mockRepo.On("CreateProduct", mock.Anything, expectedProduct).Return(returnedProduct, nil)

	result, err := svc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.SalePrice)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_RepoError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	req := productmodel.CreateProductRequest{
		Name:  "Failing Product",
		Price: 10.0,
	}

	expectedProduct := productmodel.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	mockRepo.On("CreateProduct", mock.Anything, expectedProduct).
		Return(nil, errors.New("db connection failed"))

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	name := "Updated Name"
	req := productmodel.UpdateProductRequest{
		Name: &name,
	}

	mockRepo.On("UpdateProduct", mock.Anything, "product-123", req).Return(nil)

	err := svc.UpdateProduct(context.Background(), "product-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_RepoError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	req := productmodel.UpdateProductRequest{}

	mockRepo.On("UpdateProduct", mock.Anything, "bad-id", req).
		Return(errors.New("record not found"))

	err := svc.UpdateProduct(context.Background(), "bad-id", req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NullableFieldSetToNull_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := productservice.NewProductService(mockRepo)

	req := productmodel.UpdateProductRequest{
		Description: productmodel.FieldUpdate[string]{Present: true, Value: nil},
	}

	mockRepo.On("UpdateProduct", mock.Anything, "product-123", req).Return(nil)

	err := svc.UpdateProduct(context.Background(), "product-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
