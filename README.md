# primo-test

A Go + Fiber REST API for managing products, with GORM/PostgreSQL, Swagger docs, and a full test suite (unit, integration, and E2E).

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- [swag CLI](https://github.com/swaggo/swag) (for generating Swagger docs)

## 1. Clone and install dependencies

```bash
git clone https://github.com/HorNokHuTaTo/primo-test.git
cd primo-test
go mod tidy
```

## 2. Set up environment variables

Copy the example file and fill in your own values if needed:

```bash
cp .env.example .env
```

`.env` should contain:

```env
DATABASE_URL=postgres://root:1234@localhost:5432/exam?sslmode=disable
TEST_DATABASE_URL=postgres://root:1234@localhost:5432/test_db?sslmode=disable
PORT=3000
```

## 3. Start PostgreSQL with Docker

```bash
docker-compose up -d
```

This starts a Postgres container (`my-postgres-root`) with:
- user: `root`
- password: `1234`
- default database: `exam`

Confirm it's running:

```bash
docker ps
```

## 4. Create the test database (for E2E tests)

The main `exam` database is created automatically by Docker Compose. The E2E test suite uses a **separate** database so tests don't interfere with your dev data.

```bash
docker exec -it my-postgres-root psql -U root -d exam -c "CREATE DATABASE test_db;"
```

## 5. Install the Swagger CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`, then confirm:

```bash
swag --version
```

## 6. Generate Swagger docs

```bash
swag init -g app/app.go -o docs --dir .,./cmd,./app,./config,./feature/handler,./feature/service,./feature/repository,./feature/model --parseDependency --parseInternal
```

Re-run this command any time handler annotations or request/response models change.

## 7. Run the app

```bash
go run ./cmd
```

The API will be available at `http://localhost:3000`, with interactive Swagger docs at:

```
http://localhost:3000/api-docs/index.html
```

## Running tests

Run everything (unit + integration + E2E):

```bash
go test ./... -v
```

Run a specific package only:

```bash
go test ./feature/handler/... -v
go test ./feature/service/... -v
go test ./feature/repository/... -v
```

### E2E tests

E2E tests require `TEST_DATABASE_URL` to be set and pointing at the `test_db` database created in step 4. If it's already in your `.env`, tests will pick it up automatically; otherwise set it inline:

```bash
TEST_DATABASE_URL="postgres://root:1234@localhost:5432/test_db?sslmode=disable" go test ./... -v -run TestE2E
```

If `TEST_DATABASE_URL` isn't set, E2E tests are skipped automatically rather than failing.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/product` | Create a product |
| PATCH | `/product/{id}` | Partially update a product (`description` and `sale_price` can be set to `null` to clear them) |

Full request/response schemas are available in Swagger at `/api-docs/index.html` once the app is running.
