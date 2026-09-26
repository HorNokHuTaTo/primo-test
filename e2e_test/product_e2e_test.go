package product_e2e_test

import (
	"encoding/json"
	"exam/app"
	"exam/database"
	productmodel "exam/feature/model"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB) {
	godotenv.Load("../.env")
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db := database.ConnectDB(dsn)
	db.AutoMigrate(&productmodel.Product{})

	appMock := fiber.New()
	app.SetupRoutes(appMock, db)

	return appMock, db
}

func TestE2E_CreateAndUpdateProduct(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Exec("TRUNCATE TABLE products")

	// Create
	createBody := `{"name": "Test Product", "price": 10.5}`
	createReq := httptest.NewRequest("POST", "/product", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := app.Test(createReq)
	assert.Equal(t, fiber.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&created)
	productID := created["data"].(map[string]interface{})["id"].(string)

	// Update
	updateBody := `{"price": 20.0}`
	updateReq := httptest.NewRequest("PATCH", "/product/"+productID, strings.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, _ := app.Test(updateReq)
	assert.Equal(t, fiber.StatusOK, updateResp.StatusCode)

	// Verify in DB directly
	var product productmodel.Product
	db.First(&product, "id = ?", productID)
	assert.Equal(t, 20.0, product.Price)
}
