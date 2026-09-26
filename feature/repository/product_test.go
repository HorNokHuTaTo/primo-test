package productrepository_test

import (
	"context"
	"errors"
	productmodel "exam/feature/model"
	productrepository "exam/feature/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateProduct_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	desc := "A test description"
	salePrice := 79.99
	input := productmodel.Product{
		Name:        "Test Product",
		Description: &desc,
		SalePrice:   &salePrice,
		Price:       99.99,
	}

	fakeID := "11111111-1111-1111-1111-111111111111"
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(fakeID, now, now))
	mock.ExpectCommit()

	result, err := repo.CreateProduct(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, fakeID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_MinimalFields_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	input := productmodel.Product{
		Name:  "Minimal Product",
		Price: 10.0,
	}

	fakeID := "22222222-2222-2222-2222-222222222222"
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(fakeID, now, now))
	mock.ExpectCommit()

	result, err := repo.CreateProduct(context.Background(), input)

	assert.NoError(t, err)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.SalePrice)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	input := productmodel.Product{
		Name:  "Failing Product",
		Price: 10.0,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "products"`).
		WillReturnError(errors.New("constraint violation"))
	mock.ExpectRollback()

	result, err := repo.CreateProduct(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	name := "New Name"
	req := productmodel.UpdateProductRequest{
		Name: &name,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_NullableFields_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	desc := "New description"
	req := productmodel.UpdateProductRequest{
		Description: productmodel.FieldUpdate[string]{
			Present: true,
			Value:   &desc,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_NoFieldsSet_NoOp(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	req := productmodel.UpdateProductRequest{} // nothing set

	// no Begin/Exec/Commit expected — repo should return nil without touching DB
	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_AllFields_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	name := "New Name"
	price := 99.99
	desc := "New Desc"
	salePrice := 79.99

	req := productmodel.UpdateProductRequest{
		Name:        &name,
		Price:       &price,
		Description: productmodel.FieldUpdate[string]{Present: true, Value: &desc},
		SalePrice:   productmodel.FieldUpdate[float64]{Present: true, Value: &salePrice},
	}

	mock.ExpectBegin() // struct update (name/price)
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin() // map update (description/sale_price)
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_SalePriceOnly_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	salePrice := 49.99
	req := productmodel.UpdateProductRequest{
		SalePrice: productmodel.FieldUpdate[float64]{Present: true, Value: &salePrice},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_PriceOnly_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	price := 199.99
	req := productmodel.UpdateProductRequest{
		Price: &price,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_DescriptionSetToNull_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	req := productmodel.UpdateProductRequest{
		Description: productmodel.FieldUpdate[string]{Present: true, Value: nil}, // explicit null
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_SalePriceSetToNull_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	req := productmodel.UpdateProductRequest{
		SalePrice: productmodel.FieldUpdate[float64]{Present: true, Value: nil},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_StructUpdate_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	name := "New Name"
	req := productmodel.UpdateProductRequest{Name: &name}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnError(errors.New("db connection lost"))
	mock.ExpectRollback() // GORM rolls back on error instead of committing

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_MapUpdate_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := productrepository.NewProductRepository(gormDB)

	desc := "desc"
	req := productmodel.UpdateProductRequest{
		Description: productmodel.FieldUpdate[string]{Present: true, Value: &desc},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET`).WillReturnError(errors.New("constraint violation"))
	mock.ExpectRollback()

	err := repo.UpdateProduct(context.Background(), "some-uuid", req)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
