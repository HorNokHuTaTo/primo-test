package producthandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	producthandler "exam/feature/handler"
	productmodel "exam/feature/model"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req productmodel.CreateProductRequest) (*productmodel.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productmodel.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, productID string, req productmodel.UpdateProductRequest) error {
	args := m.Called(ctx, productID, req)
	return args.Error(0)
}

func setupApp(h producthandler.IProductHandler) *fiber.App {
	app := fiber.New()
	app.Post("/product", h.CreateProduct)
	app.Patch("/product/:id", h.UpdateProduct)
	return app
}

func sendRequest(app *fiber.App, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	return app.Test(req)
}

func TestCreateProductHandler_Success(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{
		"name":  "Test Product",
		"price": 99.99,
	}

	mockService.On("CreateProduct", mock.Anything, mock.Anything).
		Return(&productmodel.Product{ID: "abc-123", Name: "Test Product", Price: 99.99}, nil)

	resp, err := sendRequest(app, "POST", "/product", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreateProductHandler_MissingRequiredField(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{
		"price": 99.99, // name missing — required field
	}

	resp, err := sendRequest(app, "POST", "/product", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockService.AssertNotCalled(t, "CreateProduct") // service should never be reached
}

func TestCreateProductHandler_InvalidBody(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	req := httptest.NewRequest("POST", "/product", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockService.AssertNotCalled(t, "CreateProduct")
}

func TestCreateProductHandler_ServiceError(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{
		"name":  "Test Product",
		"price": 99.99,
	}

	mockService.On("CreateProduct", mock.Anything, mock.Anything).
		Return(nil, errors.New("db error"))

	resp, err := sendRequest(app, "POST", "/product", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateProductHandler_Success(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{"name": "Updated Name"}

	mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := sendRequest(app, "PATCH", "/product/11111111-1111-1111-1111-111111111111", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateProductHandler_InvalidUUID(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{"name": "x"}

	resp, err := sendRequest(app, "PATCH", "/product/not-a-uuid", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockService.AssertNotCalled(t, "UpdateProduct")
}

func TestUpdateProductHandler_NullableFieldSetToNull(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{"description": nil} // explicit null

	mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.MatchedBy(func(req productmodel.UpdateProductRequest) bool {
		return req.Description.Present && req.Description.Value == nil
	})).Return(nil)

	resp, err := sendRequest(app, "PATCH", "/product/11111111-1111-1111-1111-111111111111", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateProductHandler_EmptyBody_NoFieldsSet(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{} // nothing sent

	mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := sendRequest(app, "PATCH", "/product/11111111-1111-1111-1111-111111111111", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateProductHandler_ServiceError(t *testing.T) {
	mockService := new(MockProductService)
	h := producthandler.NewProductHandler(mockService)
	app := setupApp(h)

	body := map[string]interface{}{"name": "x"}

	mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("update failed"))

	resp, err := sendRequest(app, "PATCH", "/product/11111111-1111-1111-1111-111111111111", body)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
